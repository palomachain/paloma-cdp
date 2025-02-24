package main

import (
	"log/slog"
	"os"

	"github.com/palomachain/paloma-cdp/internal/app/ingest"
	"github.com/palomachain/paloma-cdp/internal/pkg/service"
)

var version = service.DefaultVersion()

func main() {
	os.Setenv("CDP_PALOMA_RPC_ADDRESS", "https://rpc.palomachain.com:443")

	os.Setenv("CDP_PSQL_ADDRESS", "localhost:5432")
	os.Setenv("CDP_PSQL_USER", "cdp")
	os.Setenv("CDP_PSQL_PASSWORD", "trustno1")
	os.Setenv("CDP_PSQL_DATABASE", "cdp")

	svc := service.New[ingest.Configuration]().
		WithName("cdp-ingest").
		WithVersion(version).
		WithDatabase()

	if err := svc.RunWithPersistence(ingest.Run); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}
