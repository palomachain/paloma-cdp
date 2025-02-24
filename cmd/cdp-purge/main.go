package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/palomachain/paloma-cdp/internal/pkg/service"
)

var version = service.DefaultVersion()

func main() {
	os.Setenv("CDP_PSQL_ADDRESS", "localhost:5432")
	os.Setenv("CDP_PSQL_USER", "cdp")
	os.Setenv("CDP_PSQL_PASSWORD", "trustno1")
	os.Setenv("CDP_PSQL_DATABASE", "cdp")

	svc := service.New[struct{}]().
		WithName("cdp-purge").
		WithVersion(version).
		WithDatabase()

	if err := svc.RunWithPersistence(run); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context, v service.Version, db *persistence.Database, _ *struct{}) error {
	threshold := time.Now().Add(-time.Hour * 24 * 7)
	slog.Default().InfoContext(ctx, "Purging stale data.", "threshold", threshold, "version", v)

	r, err := db.NewDelete().Model(&model.PriceData{}).Where("time < ?", threshold).Exec(ctx)
	if err != nil {
		return err
	}

	count, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("data seems purged, but affected rows failed to report: %w", err)
	}

	slog.Default().InfoContext(ctx, "Data purged.", "count", count)
	return nil
}
