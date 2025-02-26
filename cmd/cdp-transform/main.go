package main

import (
	"log/slog"
	"os"

	"github.com/palomachain/paloma-cdp/internal/app/transform"
	"github.com/palomachain/paloma-cdp/internal/pkg/service"
)

var version string = "dev"

func main() {
	svc := service.New[transform.Configuration]().
		WithName("cdp-transform").
		WithVersion(version).
		WithDatabase()

	if err := svc.RunWithPersistence(transform.Run); err != nil {
		slog.Default().Error(err.Error())
		os.Exit(1)
	}
}
