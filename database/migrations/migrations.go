package migrations

import (
	"context"
	"embed"

	"github.com/nanoteck137/tunebook/database"
	"github.com/nanoteck137/tunebook/database/dup"
)

//go:embed *.sql
var migrations embed.FS

func RunMigrateUp(ctx context.Context, db *database.Database) error {
	ms, err := dup.MigrationsFromFS(migrations, ".")
	if err != nil {
		return err
	}

	return dup.RunUp(ctx, db, ms)
}
