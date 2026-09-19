// Command difftool compares user_track_stats between two tunebook databases
// and reports what was lost / changed in a rebuild.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/nanoteck137/tunebook/database"
)

type key struct {
	Type string
	Year int
	Val  int
}

type cell struct {
	Play int
	Skip int
	Time int
	Comp int
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: difftool OLD_DB NEW_DB\n\nCompares user_track_stats between two databases.\n")
	}
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}
	oldPath, newPath := flag.Arg(0), flag.Arg(1)

	ctx := context.Background()
	old, err := database.Open(oldPath)
	if err != nil {
		panic(err)
	}
	defer old.Close()
	nw, err := database.Open(newPath)
	if err != nil {
		panic(err)
	}
	defer nw.Close()

	if err := diff(ctx, old, nw); err != nil {
		panic(err)
	}
}

func load(ctx context.Context, db *database.Database) (map[string]cell, map[key]map[string]cell) {
	var rows []struct {
		Track string `db:"track_id"`
		Type  string `db:"period_type"`
		Year  int    `db:"year"`
		Val   int    `db:"period_value"`
		Play  int    `db:"play_count"`
		Skip  int    `db:"skip_count"`
		Time  int    `db:"play_time"`
		Comp  int    `db:"completion_sum"`
	}
	if err := db.Multiple(ctx, database.RawQuery{
		Query: `SELECT track_id, period_type, COALESCE(year,0) year, COALESCE(period_value,0) period_value,
		       play_count, skip_count, play_time, completion_sum FROM user_track_stats`,
	}, &rows); err != nil {
		panic(err)
	}

	total := map[string]cell{}
	byKey := map[key]map[string]cell{}
	for _, r := range rows {
		c := cell{r.Play, r.Skip, r.Time, r.Comp}
		t := total[r.Type]
		t.Play += r.Play
		t.Skip += r.Skip
		t.Time += r.Time
		t.Comp += r.Comp
		total[r.Type] = t

		k := key{r.Type, r.Year, r.Val}
		if byKey[k] == nil {
			byKey[k] = map[string]cell{}
		}
		byKey[k][r.Track] = c
	}
	return total, byKey
}

func diff(ctx context.Context, old, nw *database.Database) error {
	ot, ok := load(ctx, old)
	nt, nk := load(ctx, nw)

	fmt.Println("== Totals by period_type (old -> new) ==")
	types := []string{"all", "year", "quarter", "month"}
	for _, t := range types {
		a, b := ot[t], nt[t]
		fmt.Printf("  %-7s plays %6d -> %6d (%+d) | skip %5d -> %5d (%+d) | time %d -> %d (%+d) | comp %d -> %d\n",
			t, a.Play, b.Play, b.Play-a.Play,
			a.Skip, b.Skip, b.Skip-a.Skip,
			a.Time, b.Time, b.Time-a.Time,
			a.Comp, b.Comp)
	}

	fmt.Println("\n== Rows present in old but MISSING in new ==")
	missing := 0
	for k, m := range ok {
		n := nk[k]
		if n == nil {
			missing += len(m)
			fmt.Printf("  bucket %s/%d/%d: entire bucket missing (%d tracks)!\n", k.Type, k.Year, k.Val, len(m))
			continue
		}
		for track, c := range m {
			if _, ok2 := n[track]; !ok2 {
				missing++
				fmt.Printf("  %s year=%d val=%d track=%s plays=%d skip=%d time=%d\n",
					k.Type, k.Year, k.Val, track, c.Play, c.Skip, c.Time)
			}
		}
	}
	if missing == 0 {
		fmt.Println("  (none)")
	}

	fmt.Println("\n== Rows present in new but not in old ==")
	added := 0
	for k, n := range nk {
		m := ok[k]
		if m == nil {
			added += len(n)
			fmt.Printf("  bucket %s/%d/%d: entire bucket added (%d tracks)!\n", k.Type, k.Year, k.Val, len(n))
			continue
		}
		for track := range n {
			if _, ok2 := m[track]; !ok2 {
				added++
			}
		}
	}
	if added == 0 {
		fmt.Println("  (none)")
	}

	fmt.Println("\n== Same-key rows whose play/skip/time VALUES changed ==")
	changed := 0
	var lines []string
	for k, m := range ok {
		n := nk[k]
		if n == nil {
			continue
		}
		for track, c := range m {
			d, ok2 := n[track]
			if !ok2 {
				continue
			}
			if c.Play != d.Play || c.Skip != d.Skip || c.Time != d.Time || c.Comp != d.Comp {
				changed++
				lines = append(lines, fmt.Sprintf(
					"  %s year=%d val=%d track=%s old(play=%d skip=%d time=%d comp=%d) new(play=%d skip=%d time=%d comp=%d)",
					k.Type, k.Year, k.Val, track,
					c.Play, c.Skip, c.Time, c.Comp, d.Play, d.Skip, d.Time, d.Comp))
			}
		}
	}
	sort.Strings(lines)
	for _, l := range lines {
		fmt.Println(l)
	}
	if changed == 0 {
		fmt.Println("  (none)")
	}
	fmt.Println("\nsummary: missing rows:", missing, " added rows:", added, " value-changed rows:", changed)
	return nil
}