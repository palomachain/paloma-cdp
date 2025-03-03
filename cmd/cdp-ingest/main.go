package main

import (
	"log/slog"
	"os"

	"github.com/palomachain/paloma-cdp/internal/app/ingest"
	"github.com/palomachain/paloma-cdp/internal/pkg/service"
)

var version string = "dev"

func main() {
	svc := service.New[ingest.Configuration]().
		WithName("cdp-ingest").
		WithHealthprobe().
		WithVersion(version).
		WithDatabase()

	if err := svc.RunWithPersistence(ingest.Run); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}
