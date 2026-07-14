package main

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/dev"
	"github.com/nanoteck137/tunebook/tools/query"
	"github.com/nanoteck137/tunebook/tools/query/lexer"
	"github.com/nanoteck137/tunebook/tools/query/parser"
	"github.com/nanoteck137/tunebook/tools/query/planner"
	"github.com/nanoteck137/tunebook/tools/query/schema"
	"github.com/nanoteck137/tunebook/tools/query/sort"
	querysql "github.com/nanoteck137/tunebook/tools/query/sql"
)

func main() {
	dbTesting()

	return
	fmt.Println("=== Query Package Test ===\n")

	testLexer()
	testParser()
	testPlanner()
	testSQLCompiler()
	testSort()
	testFullPipeline()
	testDatabaseIntegration()
	testErrorMessages()

	fmt.Println("\n=== All tests completed ===")
}

func dbTesting() {
	db, err := database.Open("work/data.db")
	if err != nil {
		fmt.Printf("ERROR opening database: %v\n\n", err)
		return
	}
	defer db.Close()

	fmt.Println("Connected to work/data.db\n")

	playlistId := "nuf8abaigiryew64"

	// query := dialect.From("playlist_items").
	// 	Select("tracks.*", "playlist_items.position").
	// 	Join(
	// 		tracks.As("tracks"),
	// 		goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
	// 	)

	q := database.TrackQuery().
		SelectAppend(
			goqu.I("playlist_items.position").As("position"),
		).
		Join(
			goqu.I("playlist_items"),
			goqu.On(goqu.I("playlist_items.track_id").Eq(goqu.I("tracks.id"))),
		).Order(goqu.I("playlist_items.position").Desc())

	q = q.Where(goqu.I("playlist_items.playlist_id").Eq(playlistId))

	q, err = database.ApplyQuery(q, database.TrackSchema(), database.QueryParams{
		Filter: "tags has \"metal\"",
		Sort:   "",
	})
	if err != nil {
		fmt.Printf("ERROR apply query: %v\n\n", err)
		return
	}

	sql, params, _ := q.ToSQL()
	fmt.Printf("sql: %v\n", sql)
	fmt.Printf("params: %v\n", params)

	fmt.Printf("%v\n", database.DebugSQL(q))

	tracks, err := database.Multiple[database.PlaylistItemTrack](db, context.Background(), q)
	if err != nil {
		fmt.Printf("ERROR getting tracks: %v\n\n", err)
		return
	}

	fmt.Printf("len(tracks): %v\n", len(tracks))

	dev.Println(tracks[0])
}

func testLexer() {
	fmt.Println("--- Lexer ---\n")

	inputs := []string{
		`tags contains "rock" and year >= 1970`,
		`genre = "rock" or genre = "pop"`,
		`deleted_at is not null`,
		`genre in ("rock", "pop", "jazz")`,
		`not (year < 1970 and genre = "country")`,
		`rating >= 3.5 and active = true`,
	}

	for _, input := range inputs {
		fmt.Printf("Input:  %s\n", input)

		l := lexer.New(input)
		tokens, err := l.Scan()
		if err != nil {
			fmt.Printf("  ERROR: %v\n\n", err)
			continue
		}

		for _, tok := range tokens {
			fmt.Printf("  %s\n", tok)
		}
		fmt.Println()
	}
}

func testParser() {
	fmt.Println("--- Parser ---\n")

	inputs := []string{
		`tags contains "rock" and year >= 1970`,
		`genre = "rock" or genre = "pop"`,
		`deleted_at is not null`,
		`genre in ("rock", "pop", "jazz")`,
		`genre not in ("country", "folk")`,
		`not (year < 1970 and genre = "country")`,
		`(genre = "rock" or genre = "pop") and year >= 1970 and deleted_at is null`,
		`rating >= 3.5 and active = true`,
		`not not genre = "rock"`,
		`a = 1 or b = 2 and c = 3`,
	}

	for _, input := range inputs {
		fmt.Printf("Input:  %s\n", input)

		p := parser.New(input)
		expr, err := p.Parse()
		if err != nil {
			fmt.Printf("  PARSE ERROR: %v\n\n", err)
			continue
		}

		fmt.Printf("  AST: %s\n\n", expr)
	}
}

func testParserErrors() {
	fmt.Println("--- Parser Errors ---\n")

	inputs := []string{
		``,
		`)`,
		`year >=`,
		`= 1970`,
		`(year >= 1970`,
		`year >= 1970)`,
		`not`,
		`genre not`,
		`genre in ("rock"`,
	}

	for _, input := range inputs {
		fmt.Printf("Input: %q\n", input)

		p := parser.New(input)
		_, err := p.Parse()
		if err != nil {
			fmt.Printf("  ERROR: %v\n\n", err)
		} else {
			fmt.Printf("  (no error)\n\n")
		}
	}
}

func testPlanner() {
	fmt.Println("--- Planner ---\n")

	s := schema.New().
		AddField("name", query.TypeString).
		AddField("genre", query.TypeString).
		AddField("title", query.TypeString).
		AddField("description", query.TypeString).
		AddField("year", query.TypeInt).
		AddField("rating", query.TypeFloat).
		AddField("duration", query.TypeInt).
		AddField("active", query.TypeBool).
		AddField("deleted_at", query.TypeString, schema.Nullable())

	pl := planner.New(s)

	inputs := []string{
		`genre = "rock"`,
		`year >= 1970`,
		`rating >= 3.5`,
		`description contains "classic"`,
		`title like "%love%"`,
		`genre = "rock" and year >= 1970`,
		`genre = "rock" or genre = "pop"`,
		`not genre = "rock"`,
		`deleted_at is null`,
		`deleted_at is not null`,
		`genre in ("rock", "pop", "jazz")`,
		`genre not in ("country", "folk")`,
		`active = true`,
		`(genre = "rock" or genre = "pop") and year >= 1970 and deleted_at is not null`,
	}

	for _, input := range inputs {
		fmt.Printf("Input:  %s\n", input)

		p := parser.New(input)
		expr, err := p.Parse()
		if err != nil {
			fmt.Printf("  PARSE ERROR: %v\n\n", err)
			continue
		}

		plan, err := pl.Plan(expr)
		if err != nil {
			fmt.Printf("  PLAN ERROR: %v\n\n", err)
			continue
		}

		fmt.Printf("  Plan: %s\n\n", formatPlan(plan))
	}

	fmt.Println("--- Planner Errors ---\n")

	errorInputs := []string{
		`unknown = "value"`,
		`year contains "rock"`,
		`year like "%1970%"`,
		`name > "rock"`,
		`year = "rock"`,
		`name = 1970`,
		`name is null`,
		`year in ("rock", "pop")`,
	}

	for _, input := range errorInputs {
		fmt.Printf("Input: %s\n", input)

		p := parser.New(input)
		expr, err := p.Parse()
		if err != nil {
			fmt.Printf("  PARSE ERROR: %v\n\n", err)
			continue
		}

		_, err = pl.Plan(expr)
		if err != nil {
			fmt.Printf("  ERROR: %v\n\n", err)
		} else {
			fmt.Printf("  (no error)\n\n")
		}
	}
}

func formatPlan(plan *query.Plan) string {
	if plan.Filter == nil {
		return "(empty)"
	}
	return formatFilterNode(plan.Filter)
}

func formatFilterNode(node query.FilterNode) string {
	switch n := node.(type) {
	case *query.ComparisonNode:
		return fmt.Sprintf("Comparison(%s %s %v)", n.Field.Name, opString(n.Operator), n.Value.Value)
	case *query.ContainsNode:
		return fmt.Sprintf("Contains(%s, %v)", n.Field.Name, n.Value.Value)
	case *query.IsNullNode:
		if n.Not {
			return fmt.Sprintf("IsNotNull(%s)", n.Field.Name)
		}
		return fmt.Sprintf("IsNull(%s)", n.Field.Name)
	case *query.InNode:
		vals := make([]string, len(n.Values))
		for i, v := range n.Values {
			vals[i] = fmt.Sprintf("%v", v.Value)
		}
		not := ""
		if n.Not {
			not = "not "
		}
		return fmt.Sprintf("%sIn(%s, [%s])", not, n.Field.Name, joinStrings(vals))
	case *query.AndNode:
		return fmt.Sprintf("And(%s, %s)", formatFilterNode(n.Left), formatFilterNode(n.Right))
	case *query.OrNode:
		return fmt.Sprintf("Or(%s, %s)", formatFilterNode(n.Left), formatFilterNode(n.Right))
	case *query.NotNode:
		return fmt.Sprintf("Not(%s)", formatFilterNode(n.Expr))
	default:
		return fmt.Sprintf("%T", node)
	}
}

func opString(op query.Operator) string {
	switch op {
	case query.OpEqual:
		return "="
	case query.OpNotEqual:
		return "!="
	case query.OpGreater:
		return ">"
	case query.OpGreaterEqual:
		return ">="
	case query.OpLess:
		return "<"
	case query.OpLessEqual:
		return "<="
	case query.OpLike:
		return "like"
	default:
		return "?"
	}
}

func joinStrings(ss []string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

func testSQLCompiler() {
	fmt.Println("--- SQL Compiler ---\n")

	compiler := querysql.NewCompiler()

	testComparison(compiler)
	testLogicalOperators(compiler)
	testSpecialOperators(compiler)
	testOrdering(compiler)
	testComplexQuery(compiler)
}

func testComparison(c *querysql.Compiler) {
	fmt.Println("  -- Comparison Operators --\n")

	tests := []struct {
		name string
		plan *query.Plan
	}{
		{
			name: "Equal",
			plan: &query.Plan{
				Filter: &query.ComparisonNode{
					Field:    &query.Field{Name: "name", Type: query.TypeString},
					Operator: query.OpEqual,
					Value:    query.Value{Type: query.TypeString, Value: "rock"},
				},
			},
		},
		{
			name: "Not Equal",
			plan: &query.Plan{
				Filter: &query.ComparisonNode{
					Field:    &query.Field{Name: "genre", Type: query.TypeString},
					Operator: query.OpNotEqual,
					Value:    query.Value{Type: query.TypeString, Value: "pop"},
				},
			},
		},
		{
			name: "Greater Equal",
			plan: &query.Plan{
				Filter: &query.ComparisonNode{
					Field:    &query.Field{Name: "year", Type: query.TypeInt},
					Operator: query.OpGreaterEqual,
					Value:    query.Value{Type: query.TypeInt, Value: 1970},
				},
			},
		},
		{
			name: "Less Than",
			plan: &query.Plan{
				Filter: &query.ComparisonNode{
					Field:    &query.Field{Name: "rating", Type: query.TypeFloat},
					Operator: query.OpLess,
					Value:    query.Value{Type: query.TypeFloat, Value: 3.5},
				},
			},
		},
		{
			name: "Like",
			plan: &query.Plan{
				Filter: &query.ComparisonNode{
					Field:    &query.Field{Name: "title", Type: query.TypeString},
					Operator: query.OpLike,
					Value:    query.Value{Type: query.TypeString, Value: "%love%"},
				},
			},
		},
	}

	for _, tt := range tests {
		printPlan(c, "  "+tt.name, tt.plan)
	}
}

func testLogicalOperators(c *querysql.Compiler) {
	fmt.Println("  -- Logical Operators --\n")

	tests := []struct {
		name string
		plan *query.Plan
	}{
		{
			name: "AND",
			plan: &query.Plan{
				Filter: &query.AndNode{
					Left: &query.ComparisonNode{
						Field:    &query.Field{Name: "genre", Type: query.TypeString},
						Operator: query.OpEqual,
						Value:    query.Value{Type: query.TypeString, Value: "rock"},
					},
					Right: &query.ComparisonNode{
						Field:    &query.Field{Name: "year", Type: query.TypeInt},
						Operator: query.OpGreaterEqual,
						Value:    query.Value{Type: query.TypeInt, Value: 1970},
					},
				},
			},
		},
		{
			name: "OR",
			plan: &query.Plan{
				Filter: &query.OrNode{
					Left: &query.ComparisonNode{
						Field:    &query.Field{Name: "genre", Type: query.TypeString},
						Operator: query.OpEqual,
						Value:    query.Value{Type: query.TypeString, Value: "rock"},
					},
					Right: &query.ComparisonNode{
						Field:    &query.Field{Name: "genre", Type: query.TypeString},
						Operator: query.OpEqual,
						Value:    query.Value{Type: query.TypeString, Value: "pop"},
					},
				},
			},
		},
		{
			name: "NOT",
			plan: &query.Plan{
				Filter: &query.NotNode{
					Expr: &query.ComparisonNode{
						Field:    &query.Field{Name: "genre", Type: query.TypeString},
						Operator: query.OpEqual,
						Value:    query.Value{Type: query.TypeString, Value: "country"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		printPlan(c, "  "+tt.name, tt.plan)
	}
}

func testSpecialOperators(c *querysql.Compiler) {
	fmt.Println("  -- Special Operators --\n")

	tests := []struct {
		name string
		plan *query.Plan
	}{
		{
			name: "IS NULL",
			plan: &query.Plan{
				Filter: &query.IsNullNode{
					Field: &query.Field{Name: "deleted_at", Type: query.TypeString, Nullable: true},
				},
			},
		},
		{
			name: "IS NOT NULL",
			plan: &query.Plan{
				Filter: &query.IsNullNode{
					Field: &query.Field{Name: "deleted_at", Type: query.TypeString, Nullable: true},
					Not:   true,
				},
			},
		},
		{
			name: "CONTAINS",
			plan: &query.Plan{
				Filter: &query.ContainsNode{
					Field: &query.Field{Name: "description", Type: query.TypeString},
					Value: query.Value{Type: query.TypeString, Value: "classic"},
				},
			},
		},
		{
			name: "IN",
			plan: &query.Plan{
				Filter: &query.InNode{
					Field: &query.Field{Name: "genre", Type: query.TypeString},
					Values: []query.Value{
						{Type: query.TypeString, Value: "rock"},
						{Type: query.TypeString, Value: "pop"},
						{Type: query.TypeString, Value: "jazz"},
					},
				},
			},
		},
		{
			name: "NOT IN",
			plan: &query.Plan{
				Filter: &query.InNode{
					Field: &query.Field{Name: "status", Type: query.TypeString},
					Values: []query.Value{
						{Type: query.TypeString, Value: "deleted"},
						{Type: query.TypeString, Value: "archived"},
					},
					Not: true,
				},
			},
		},
	}

	for _, tt := range tests {
		printPlan(c, "  "+tt.name, tt.plan)
	}
}

func testOrdering(c *querysql.Compiler) {
	fmt.Println("  -- Ordering --\n")

	tests := []struct {
		name string
		plan *query.Plan
	}{
		{
			name: "Single ASC",
			plan: &query.Plan{
				OrderBy: []query.Ordering{
					&query.FieldOrdering{
						Field: &query.Field{Name: "name", Type: query.TypeString},
						Dir:   query.DirAsc,
					},
				},
			},
		},
		{
			name: "Multiple",
			plan: &query.Plan{
				OrderBy: []query.Ordering{
					&query.FieldOrdering{
						Field: &query.Field{Name: "artist", Type: query.TypeString},
						Dir:   query.DirAsc,
					},
					&query.FieldOrdering{
						Field: &query.Field{Name: "year", Type: query.TypeInt},
						Dir:   query.DirDesc,
					},
				},
			},
		},
		{
			name: "Random",
			plan: &query.Plan{
				OrderBy: []query.Ordering{
					&query.RandomOrdering{},
				},
			},
		},
	}

	for _, tt := range tests {
		printPlan(c, "  "+tt.name, tt.plan)
	}
}

func testComplexQuery(c *querysql.Compiler) {
	fmt.Println("  -- Complex Query --\n")

	plan := &query.Plan{
		Filter: &query.AndNode{
			Left: &query.OrNode{
				Left: &query.ComparisonNode{
					Field:    &query.Field{Name: "genre", Type: query.TypeString},
					Operator: query.OpEqual,
					Value:    query.Value{Type: query.TypeString, Value: "rock"},
				},
				Right: &query.ComparisonNode{
					Field:    &query.Field{Name: "genre", Type: query.TypeString},
					Operator: query.OpEqual,
					Value:    query.Value{Type: query.TypeString, Value: "pop"},
				},
			},
			Right: &query.AndNode{
				Left: &query.ComparisonNode{
					Field:    &query.Field{Name: "year", Type: query.TypeInt},
					Operator: query.OpGreaterEqual,
					Value:    query.Value{Type: query.TypeInt, Value: 1970},
				},
				Right: &query.IsNullNode{
					Field: &query.Field{Name: "deleted_at", Type: query.TypeString, Nullable: true},
					Not:   true,
				},
			},
		},
		OrderBy: []query.Ordering{
			&query.FieldOrdering{
				Field: &query.Field{Name: "year", Type: query.TypeInt},
				Dir:   query.DirDesc,
			},
			&query.FieldOrdering{
				Field: &query.Field{Name: "title", Type: query.TypeString},
				Dir:   query.DirAsc,
			},
		},
	}

	printPlan(c, "  Complex Query", plan)
}

func printPlan(c *querysql.Compiler, name string, plan *query.Plan) {
	fmt.Printf("  %s\n", name)

	res, err := c.Compile(plan)
	if err != nil {
		fmt.Printf("    ERROR: %v\n\n", err)
		return
	}

	if res.Where != nil {
		sql, args, err := goqu.From("dummy").Where(res.Where).ToSQL()
		if err != nil {
			fmt.Printf("    TO SQL ERROR: %v\n\n", err)
			return
		}
		// Extract just the WHERE clause
		whereIdx := len("SELECT * FROM \"dummy\" WHERE ")
		if len(sql) > whereIdx {
			sql = sql[whereIdx:]
		}
		fmt.Printf("    WHERE: %s\n", sql)
		if len(args) > 0 {
			fmt.Printf("    ARGS:  %v\n", args)
		}
	}

	if len(res.Order) > 0 {
		for i, ord := range res.Order {
			sql, _, err := goqu.From("dummy").Order(ord).ToSQL()
			if err != nil {
				fmt.Printf("    TO SQL ERROR: %v\n\n", err)
				continue
			}
			// Extract just the ORDER BY clause
			orderIdx := len("SELECT * FROM \"dummy\" ORDER BY ")
			if len(sql) > orderIdx {
				sql = sql[orderIdx:]
			}
			if i == 0 {
				fmt.Printf("    ORDER: %s\n", sql)
			} else {
				fmt.Printf("           %s\n", sql)
			}
		}
	}

	if res.Where == nil && len(res.Order) == 0 {
		fmt.Printf("    (empty plan)\n")
	}

	fmt.Println()
}

func testFullPipeline() {
	fmt.Println("--- Full Pipeline (Input -> Lexer -> Parser -> Planner -> SQL) ---\n")

	s := schema.New().
		AddField("name", query.TypeString).
		AddField("genre", query.TypeString).
		AddField("title", query.TypeString).
		AddField("description", query.TypeString).
		AddField("year", query.TypeInt).
		AddField("rating", query.TypeFloat).
		AddField("duration", query.TypeInt).
		AddField("active", query.TypeBool).
		AddField("deleted_at", query.TypeString, schema.Nullable())

	pl := planner.New(s)
	compiler := querysql.NewCompiler()

	inputs := []string{
		`genre = "rock"`,
		`year >= 1970`,
		`rating >= 3.5`,
		`description contains "classic"`,
		`title like "%love%"`,
		`genre = "rock" and year >= 1970`,
		`genre = "rock" or genre = "pop"`,
		`not genre = "country"`,
		`deleted_at is null`,
		`deleted_at is not null`,
		`genre in ("rock", "pop", "jazz")`,
		`genre not in ("country", "folk")`,
		`active = true`,
		`(genre = "rock" or genre = "pop") and year >= 1970 and deleted_at is not null`,
	}

	for _, input := range inputs {
		fmt.Printf("Input: %s\n", input)

		p := parser.New(input)
		expr, err := p.Parse()
		if err != nil {
			fmt.Printf("  PARSE ERROR: %v\n\n", err)
			continue
		}

		plan, err := pl.Plan(expr)
		if err != nil {
			fmt.Printf("  PLAN ERROR: %v\n\n", err)
			continue
		}

		res, err := compiler.Compile(plan)
		if err != nil {
			fmt.Printf("  COMPILE ERROR: %v\n\n", err)
			continue
		}

		if res.Where != nil {
			sql, args, err := goqu.From("dummy").Where(res.Where).ToSQL()
			if err != nil {
				fmt.Printf("  TO SQL ERROR: %v\n\n", err)
				continue
			}
			// Extract just the WHERE clause
			whereIdx := len("SELECT * FROM \"dummy\" WHERE ")
			if len(sql) > whereIdx {
				sql = sql[whereIdx:]
			}
			fmt.Printf("  SQL:   %s\n", sql)
			if len(args) > 0 {
				fmt.Printf("  ARGS:  %v\n", args)
			}
		}
		fmt.Println()
	}
}

func testSort() {
	fmt.Println("--- Sort Parser ---\n")

	inputs := []string{
		"",
		"random",
		"score",
		"recent",
		"shuffle=1234",
		"+name",
		"-year",
		"+artist,-year",
		"name asc",
		"year desc",
		"artist asc, year desc",
		"name",
		"+artist, year desc, -duration",
		"  +artist  ,  -year  ",
		"name ASC, year DESC",
		"name nulls first",
		"year desc nulls last",
		"+artist nulls first, -year nulls last",
		"name asc nulls first, year desc nulls last",
	}

	for _, input := range inputs {
		fmt.Printf("Input: %q\n", input)
		s, err := sort.Parse(input)
		if err != nil {
			fmt.Printf("  ERROR: %v\n\n", err)
			continue
		}

		fmt.Printf("  Orderings: %d\n", len(s.Orderings))
		for i, o := range s.Orderings {
			switch o := o.(type) {
			case *query.FieldOrdering:
				dir := "ASC"
				if o.Dir == query.DirDesc {
					dir = "DESC"
				}
				nullStr := ""
				switch o.NullOrder {
				case query.NullOrderingFirst:
					nullStr = " NULLS FIRST"
				case query.NullOrderingLast:
					nullStr = " NULLS LAST"
				}
				fmt.Printf("    [%d] Field: %s %s%s\n", i, o.Field.Name, dir, nullStr)
			case *query.RandomOrdering:
				fmt.Printf("    [%d] Random\n", i)
			case *query.ShuffleOrdering:
				fmt.Printf("    [%d] Shuffle (seed: %d)\n", i, o.Seed)
			case *query.ScoreOrdering:
				fmt.Printf("    [%d] Score\n", i)
			}
		}
		fmt.Println()
	}

	fmt.Println("--- Sort Errors ---\n")

	errorInputs := []string{
		"shuffle=abc",
		"+",
		"-",
		" asc",
		" desc",
	}

	for _, input := range errorInputs {
		fmt.Printf("Input: %q\n", input)
		_, err := sort.Parse(input)
		if err != nil {
			fmt.Printf("  ERROR: %v\n\n", err)
		} else {
			fmt.Printf("  (no error)\n\n")
		}
	}
}

func testDatabaseIntegration() {
	fmt.Println("--- Database Integration Test ---\n")

	db, err := database.Open("work/data.db")
	if err != nil {
		fmt.Printf("ERROR opening database: %v\n\n", err)
		return
	}
	defer db.Close()

	fmt.Println("Connected to work/data.db\n")

	s := schema.New().
		AddField("id", query.TypeString, schema.Column("tracks.id")).
		AddField("name", query.TypeString, schema.Column("tracks.name")).
		AddField("number", query.TypeInt, schema.Column("tracks.number"), schema.Nullable()).
		AddField("duration", query.TypeInt, schema.Column("tracks.duration"), schema.Nullable()).
		AddField("year", query.TypeInt, schema.Column("tracks.year"), schema.Nullable()).
		AddField("albumId", query.TypeString, schema.Column("tracks.album_id")).
		AddField("artistId", query.TypeString, schema.Column("tracks.artist_id")).
		AddField("albumName", query.TypeString, schema.Column("albums.name")).
		AddField("artistName", query.TypeString, schema.Column("artists.name")).
		AddField("tags", query.TypeRelation, schema.Relation("tracks_tags", "track_id", "tag_slug", query.TypeString, "tracks.id")).
		AddField("featuringArtist", query.TypeRelation, schema.Relation("tracks_featuring_artists", "track_id", "artist_id", query.TypeString, "tracks.id")).
		AddField("ratingRelation", query.TypeRelation, schema.Relation("track_ratings", "track_id", "rating_value", query.TypeInt, "tracks.id")).
		AddField("created", query.TypeInt, schema.Column("tracks.created")).
		AddField("updated", query.TypeInt, schema.Column("tracks.updated")).
		SetDefaultSort(
			&query.FieldOrdering{
				Field: &query.Field{Name: "artistName"},
				Dir:   query.DirAsc,
			},
			&query.FieldOrdering{
				Field: &query.Field{Name: "year"},
				Dir:   query.DirDesc,
			},
		)

	pl := planner.New(s)
	compiler := querysql.NewCompiler()

	tests := []struct {
		name  string
		query string
		sort  string
	}{
		{
			name:  "Default sort (no sort provided)",
			query: `year >= 2000`,
		},
		{
			name:  "Year >= 2000",
			query: `year >= 2000`,
			sort:  "name asc",
		},
		{
			name:  "Duration between 180 and 300",
			query: `duration >= 180 and duration <= 300`,
		},
		{
			name:  "Artist name contains 'The'",
			query: `artistName contains "The"`,
		},
		{
			name:  "Year is not null",
			query: `year is not null`,
		},
		{
			name:  "Complex: year >= 2000 and duration > 120",
			query: `year >= 2000 and duration > 120`,
		},
		{
			name:  "Album name contains 'live'",
			query: `albumName contains "live"`,
		},
		{
			name:  "Tags has 'soundtrack' (relation)",
			query: `tags has "soundtrack"`,
		},
		{
			name:  "Tags has 'rock' (relation)",
			query: `tags has "rock"`,
		},
		{
			name:  "Tags NOT has 'soundtrack' (relation)",
			query: `not tags has "soundtrack"`,
		},
		{
			name:  "Tags has 'rock' and year >= 2000 (relation)",
			query: `tags has "rock" and year >= 2000`,
		},
		{
			name:  "Tags has 'live'",
			query: `tags has "live"`,
		},
		{
			name:  "Featuring artist has specific ID",
			query: `featuringArtist has "artist_123"`,
		},
		{
			name:  "Featuring artist NOT has specific ID",
			query: `not featuringArtist has "artist_123"`,
		},
		{
			name:  "Tags has 'rock' and featuring artist has ID",
			query: `tags has "rock" and featuringArtist has "vda30ry85z"`,
		},
		{
			name:  "Rating relation has integer value 5",
			query: `ratingRelation has 5`,
		},
		{
			name:  "Rating relation has integer value 4 or 5",
			query: `ratingRelation has 4 or ratingRelation has 5`,
		},
		{
			name:  "Precedence: year >= 2000 or year <= 1990 and duration > 300",
			query: `year >= 2000 or year <= 1990 and duration > 300`,
		},
		{
			name:  "Precedence with parens: (year >= 2000 or year <= 1990) and duration > 300",
			query: `(year >= 2000 or year <= 1990) and duration > 300`,
		},
		{
			name:  "Precedence: tags has 'rock' or tags has 'pop' and year >= 2000",
			query: `tags has "rock" or tags has "pop" and year >= 2000`,
		},
		{
			name:  "Precedence with parens: (tags has 'rock' or tags has 'pop') and year >= 2000",
			query: `(tags has "rock" or tags has "pop") and year >= 2000`,
		},
		{
			name:  "Complex nested: (year >= 2000 and (tags has 'rock' or tags has 'metal')) or duration < 120",
			query: `(year >= 2000 and (tags has "rock" or tags has "metal")) or duration < 120`,
		},
		{
			name:  "Deep nesting: ((year >= 2000 or year <= 1980) and (tags has 'rock' or tags has 'pop')) and not tags has 'live'",
			query: `((year >= 2000 or year <= 1980) and (tags has "rock" or tags has "pop")) and not tags has "live"`,
		},
		{
			name:  "Multiple OR with AND: year >= 2000 and (tags has 'rock' or tags has 'pop' or tags has 'metal')",
			query: `year >= 2000 and (tags has "rock" or tags has "pop" or tags has "metal")`,
		},
		{
			name:  "NOT with parens: not (tags has 'live' or tags has 'acoustic')",
			query: `not (tags has "live" or tags has "acoustic")`,
		},
		{
			name:  "NOT without parens: not tags has 'live' or not tags has 'acoustic'",
			query: `not tags has "live" or not tags has "acoustic"`,
		},
		{
			name:  "Sort by name ascending",
			query: `year >= 2000`,
			sort:  "+name",
		},
		{
			name:  "Sort by year descending, name ascending",
			query: `year >= 2000`,
			sort:  "-year,+name",
		},
		{
			name:  "Sort by duration ascending",
			query: `duration > 0`,
			sort:  "duration asc",
		},
		{
			name:  "Sort by artist name and year",
			query: `year >= 2000`,
			sort:  "artistName asc, year desc",
		},
		{
			name:  "Random sort",
			query: `year >= 2000`,
			sort:  "random",
		},
		{
			name:  "Recent sort",
			query: `year >= 2000`,
			sort:  "recent",
		},
		{
			name:  "Sort with nulls first",
			query: `year >= 2000`,
			sort:  "year desc nulls first",
		},
		{
			name:  "Sort with nulls last",
			query: `year >= 2000`,
			sort:  "year asc nulls last",
		},
		{
			name:  "Sort with mixed null ordering",
			query: `year >= 2000`,
			sort:  "artistName asc nulls first, year desc nulls last",
		},
	}

	for _, tt := range tests {
		fmt.Printf("--- %s ---\n", tt.name)
		fmt.Printf("Query: %s\n", tt.query)
		if tt.sort != "" {
			fmt.Printf("Sort:  %s\n", tt.sort)
		} else {
			fmt.Printf("Sort:  (using default)\n")
		}

		p := parser.New(tt.query)
		expr, err := p.Parse()
		if err != nil {
			fmt.Printf("  PARSE ERROR: %v\n\n", err)
			continue
		}

		plan, err := pl.Plan(expr)
		if err != nil {
			fmt.Printf("  PLAN ERROR: %v\n\n", err)
			continue
		}

		// Parse sorting if provided, otherwise use default
		var sortOrderings []query.Ordering
		if tt.sort != "" {
			sortObj, err := sort.Parse(tt.sort)
			if err != nil {
				fmt.Printf("  SORT ERROR: %v\n\n", err)
				continue
			}
			sortOrderings = sortObj.Orderings
		}

		// Resolve field names using the schema (applies default if sortOrderings is empty)
		resolvedOrderings, err := pl.ResolveSort(sortOrderings)
		if err != nil {
			fmt.Printf("  SORT RESOLUTION ERROR: %v\n\n", err)
			continue
		}
		plan.OrderBy = resolvedOrderings

		result, err := compiler.Compile(plan)
		if err != nil {
			fmt.Printf("  COMPILE ERROR: %v\n\n", err)
			continue
		}

		count, err := executeTrackQuery(db, result)
		if err != nil {
			fmt.Printf("  EXEC ERROR: %v\n\n", err)
			continue
		}

		fmt.Printf("Results: %d tracks\n\n", count)
	}
}

func executeTrackQuery(executor database.Executor, plan *querysql.CompileResult) (int, error) {
	q := database.TrackQuery().Prepared(true)

	if plan.Where != nil {
		q = q.Where(plan.Where)
	}

	if len(plan.Order) > 0 {
		q = q.Order(plan.Order...)
	}

	tracks, err := database.Multiple[database.Track](executor, context.Background(), q)
	if err != nil {
		return 0, fmt.Errorf("execute: %w", err)
	}

	sqlStr, args, err := q.ToSQL()
	if err != nil {
		return 0, fmt.Errorf("SQL generation error: %w", err)
	}

	const debugPrint = false
	if debugPrint {
		fmt.Printf("sqlStr: %v\n", sqlStr)
		fmt.Printf("args: %v\n", args)
	}

	return len(tracks), nil
}

func testErrorMessages() {
	fmt.Println("--- Error Messages ---\n")

	s := schema.New().
		AddField("id", query.TypeString, schema.Column("tracks.id")).
		AddField("name", query.TypeString, schema.Column("tracks.name")).
		AddField("year", query.TypeInt, schema.Column("tracks.year"), schema.Nullable()).
		AddField("duration", query.TypeInt, schema.Column("tracks.duration"), schema.Nullable()).
		AddField("rating", query.TypeFloat, schema.Column("tracks.rating")).
		AddField("active", query.TypeBool, schema.Column("tracks.active")).
		AddField("tag", query.TypeRelation, schema.Relation("tracks_tags", "track_id", "tag_slug", query.TypeString, "tracks.id")).
		AddField("ratingRelation", query.TypeRelation, schema.Relation("track_ratings", "track_id", "rating_value", query.TypeInt, "tracks.id"))

	pl := planner.New(s)

	fmt.Println("=== Lexer Errors ===\n")

	lexerErrors := []string{
		`name = "unterminated string`,
		`year >= 1970 @`,
		`name = "test\"`,
		`year >= 123abc`,
	}

	for _, input := range lexerErrors {
		fmt.Printf("Input: %s\n", input)
		l := lexer.New(input)
		_, err := l.Scan()
		if err != nil {
			fmt.Printf("  Lexer Error: %v\n\n", err)
		}
	}

	fmt.Println("=== Parser Errors ===\n")

	parserErrors := []string{
		`name = `,
		`= "rock"`,
		`name = "rock" and`,
		`name = "rock" or`,
		`(name = "rock"`,
		`name = "rock")`,
		`name = "rock" and and year >= 1970`,
		`year >= >= 1970`,
		`name in`,
		`name in (`,
		`name in ("rock"`,
		`name in ("rock",)`,
	}

	for _, input := range parserErrors {
		fmt.Printf("Input: %s\n", input)
		p := parser.New(input)
		_, err := p.Parse()
		if err != nil {
			fmt.Printf("  Parser Error: %v\n\n", err)
		}
	}

	fmt.Println("=== Planner Errors ===\n")

	plannerErrors := []string{
		`unknownField = "value"`,
		`year = "string on int field"`,
		`name = 1970`,
		`rating = "string on float field"`,
		`active = "string on bool field"`,
		`year > "string comparison"`,
		`name > "string greater than"`,
		`name < "string less than"`,
		`name >= "string greater equal"`,
		`name <= "string less equal"`,
		`year contains "contains on int"`,
		`rating contains "contains on float"`,
		`active contains "contains on bool"`,
		`year like "like on int"`,
		`rating like "like on float"`,
		`active like "like on bool"`,
		`name is null`,
		`name is not null`,
		`year in ("string", "values")`,
		`tag = "equality on relation"`,
		`tag > "comparison on relation"`,
		`tag contains "contains on relation"`,
		`tag has 123`,
		`ratingRelation has "string on int relation"`,
	}

	for _, input := range plannerErrors {
		fmt.Printf("Input: %s\n", input)
		p := parser.New(input)
		expr, err := p.Parse()
		if err != nil {
			fmt.Printf("  Parser Error: %v\n\n", err)
			continue
		}
		_, err = pl.Plan(expr)
		if err != nil {
			fmt.Printf("  Planner Error: %v\n\n", err)
		}
	}
}
