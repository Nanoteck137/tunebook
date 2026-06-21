package dup

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/nanoteck137/tunebook/database"
)

var dialect = goqu.Dialect("sqlite3")

type Migration struct {
	Version int64
	Name    string
	UpSQL   string
}

func parseMigrationFile(name string, content string) (Migration, error) {
	parts := strings.SplitN(name, "_", 2)
	if len(parts) < 2 {
		return Migration{}, fmt.Errorf(
			"dup: invalid migration filename %q", name)
	}

	version, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Migration{}, fmt.Errorf(
			"dup: invalid version in filename %q: %w", name, err)
	}

	displayName := strings.TrimSuffix(parts[1], ".sql")

	return Migration{
		Version: version,
		Name:    displayName,
		UpSQL:   strings.TrimSpace(content),
	}, nil
}

func MigrationsFromFS(fsys embed.FS, dir string) ([]Migration, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("dup: reading embedded dir %q: %w", dir, err)
	}

	var migrations []Migration

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		data, err := fs.ReadFile(fsys, path.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("dup: reading %q: %w", entry.Name(), err)
		}

		m, err := parseMigrationFile(entry.Name(), string(data))
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, m)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func ensureMigrationsTable(ctx context.Context, db database.Executor) error {
	_, err := db.Exec(ctx, database.RawQuery{
		Query: `CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at INTEGER NOT NULL
		)`,
	})

	return err
}

func appliedVersions(
	ctx context.Context,
	db database.Executor,
) (map[int64]struct{}, error) {
	query := dialect.From("migrations").
		Select("version").
		Order(goqu.I("version").Asc())

	versions, err := database.Multiple[int64](db, ctx, query)
	if err != nil {
		return nil, err
	}

	applied := make(map[int64]struct{}, len(versions))
	for _, v := range versions {
		applied[v] = struct{}{}
	}

	return applied, nil
}

func RunUp(
	ctx context.Context,
	db *database.Database,
	migrations []Migration,
) error {
	err := ensureMigrationsTable(ctx, db)
	if err != nil {
		return fmt.Errorf("dup: create migrations table: %w", err)
	}

	applied, err := appliedVersions(ctx, db)
	if err != nil {
		return fmt.Errorf("dup: read applied versions: %w", err)
	}

	apply := func(m Migration) error {
		slog.Info(
			"dup: applying migration",
			"version", m.Version,
			"name", m.Name,
		)

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf(
				"dup: begin tx for migration %d: %w", m.Version, err)
		}
		defer tx.Rollback()

		_, err = tx.Exec(ctx, database.RawQuery{Query: m.UpSQL})
		if err != nil {
			return fmt.Errorf(
				"dup: run migration %d (%s): %w", m.Version, m.Name, err)
		}

		now := time.Now().UnixMilli()
		insert := dialect.Insert("migrations").Rows(goqu.Record{
			"version":    m.Version,
			"name":       m.Name,
			"applied_at": now,
		})

		_, err = tx.Exec(ctx, insert)
		if err != nil {
			return fmt.Errorf("dup: record migration %d: %w", m.Version, err)
		}

		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("dup: commit migration %d: %w", m.Version, err)
		}

		return nil
	}

	for _, m := range migrations {
		if _, ok := applied[m.Version]; ok {
			slog.Info(
				"dup: skipping migration",
				"version", m.Version,
				"name", m.Name,
			)

			continue
		}

		if m.UpSQL == "" {
			slog.Info(
				"dup: skipping migration, empty up", 
				"version", m.Version, 
				"name", m.Name,
			)

			continue
		}

		err := apply(m)
		if err != nil {
			return err
		}
	}

	return nil
}
