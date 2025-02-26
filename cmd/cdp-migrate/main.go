package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/palomachain/paloma-cdp/internal/pkg/service"
)

var version string = "dev"

func main() {
	svc := service.New[struct{}]().
		WithName("cdp-migrate").
		WithVersion(version).
		WithDatabase()

	if err := svc.RunWithPersistence(run); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, v string, db *persistence.Database, _ *struct{}) error {
	slog.Default().InfoContext(ctx, "Service running.", "version", v)
	if err := db.Migrate(ctx); err != nil {
		liblog.WithError(ctx, err, "Failed to migrate database, manual intervention required!")
		return err
	}

	return nil
}
