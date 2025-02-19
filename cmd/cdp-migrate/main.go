package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/palomachain/paloma-cdp/internal/pkg/service"
)

func main() {
	os.Setenv("CDP_PSQL_ADDRESS", "localhost:5432")
	os.Setenv("CDP_PSQL_USER", "cdp")
	os.Setenv("CDP_PSQL_PASSWORD", "trustno1")
	os.Setenv("CDP_PSQL_DATABASE", "cdp")

	svc := service.New[struct{}]().
		WithName("cdp-migrate").
		WithDatabase()

	if err := svc.RunWithPersistence(run); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, db *persistence.Database, _ *struct{}) error {
	if err := db.Migrate(ctx); err != nil {
		liblog.WithError(ctx, err, "Failed to migrate database, manual intervention required!")
		return err
	}

	return nil
}
