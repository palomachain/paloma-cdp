package main

import (
	"log/slog"
	"os"

	"github.com/palomachain/paloma-cdp/internal/app/rest"
	"github.com/palomachain/paloma-cdp/internal/pkg/service"
)

var version = service.DefaultVersion()

func main() {
	os.Setenv("CDP_PSQL_ADDRESS", "localhost:5432")
	os.Setenv("CDP_PSQL_USER", "cdp")
	os.Setenv("CDP_PSQL_PASSWORD", "trustno1")
	os.Setenv("CDP_PSQL_DATABASE", "cdp")

	svc := service.New[rest.Configuration]().
		WithName("cdp-rest").
		WithVersion(version).
		WithDatabase()

	if err := svc.RunWithPersistence(rest.Run); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}
