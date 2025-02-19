package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	v1 "github.com/palomachain/paloma-cdp/internal/app/rest/v1"
	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/response/gzip"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
)

type Configuration struct {
	HttpPort string `env:"CDP_REST_API_PORT" envDefault:"8011"`
	HttpHost string `env:"CDP_REST_API_HOST" envDefault:"localhost"`
}

func Run(
	ctx context.Context,
	db *persistence.Database,
	cfg *Configuration,
) error {
	s := web.NewService(openapi31.NewReflector())

	// Init API documentation schema.
	s.OpenAPISchema().SetTitle("Basic Example")
	s.OpenAPISchema().SetDescription("This app showcases a trivial REST API.")
	s.OpenAPISchema().SetVersion("v1.2.3")

	s.Wrap(
		gzip.Middleware,
	)

	s.Use(
		liblog.Middleware("cdp-rest"),
	)

	// Add use case handler to router.
	s.Get("/api/v1/symbol/{name}", v1.SymbolInteractor(ctx, db))
	s.Get("/api/v1/symbol/{name}/bars", v1.BarsInteractor(db))
	s.Get("/api/v1/symbols", v1.SymbolsInteractor(ctx, db))

	// Swagger UI endpoint at /docs.
	s.Docs("/docs", swgui.New)

	// Start server.
	binding := fmt.Sprintf("%s:%s", cfg.HttpHost, cfg.HttpPort)
	srv := http.Server{
		Addr:    binding,
		Handler: s,
	}
	slog.Default().InfoContext(ctx, "Service running.", "binding", binding)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				liblog.WithError(ctx, err, "HTTP server failed.")
				panic(err)
			}
		}
	}()

	<-ctx.Done()
	srv.Shutdown(ctx)
	slog.Default().InfoContext(ctx, "Service stopped.")

	return nil
}
