package main

import (
	"log/slog"
	"os"

	"github.com/palomachain/paloma-cdp/internal/app/rest"
	"github.com/palomachain/paloma-cdp/internal/pkg/service"
)

var version string = "dev"

func main() {
	svc := service.New[rest.Configuration]().
		WithName("cdp-rest").
		WithVersion(version).
		WithDatabase()

	if err := svc.RunWithPersistence(rest.Run); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}
