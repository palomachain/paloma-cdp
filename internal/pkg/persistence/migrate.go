package persistence

import (
	"context"
	"embed"
	"fmt"
	"log/slog"

	"github.com/uptrace/bun/migrate"
)

//go:embed sql/*.sql
var sqlMigrations embed.FS

func (db *Database) Migrate(ctx context.Context) (err error) {
	logger := slog.Default().With("component", "migrator")
	Migrations := migrate.NewMigrations()
	if err := Migrations.Discover(sqlMigrations); err != nil {
		return fmt.Errorf("failed to discover migrations: %w", err)
	}

	m := migrate.NewMigrator(db.DB, Migrations)
	err = m.Init(ctx)
	if err != nil {
		return fmt.Errorf("failed to init migrator: %w", err)
	}

	err = m.Lock(ctx)
	if err != nil {
		return fmt.Errorf("failed to lock migrator: %w", err)
	}

	defer func() {
		if err := m.Unlock(ctx); err != nil {
			panic(err)
		}
	}()

	g, err := m.Migrate(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to run migrations. Attempting rollback", "error", err)
		if _, err := m.Rollback(ctx); err != nil {
			panic(err)
		}
		return err
	}

	for _, v := range g.Migrations.Applied() {
		logger.InfoContext(ctx, "Applied migration", "id", v.ID, "name", v.Name)
	}

	return nil
}
