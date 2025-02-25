package rest

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
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
	v string,
	db *persistence.Database,
	cfg *Configuration,
) error {
	s := web.NewService(openapi31.NewReflector())

	s.OpenAPISchema().SetTitle("Paloma Chain Data Provider - REST API")
	s.OpenAPISchema().SetDescription("This API grants access to live and historic chain data from Paloma. The initial feature set was built to satisfy charting solutions, but may be extended in the future.")
	s.OpenAPISchema().SetVersion(v)

	s.Wrap(
		gzip.Middleware,
	)

	s.Use(
		liblog.Middleware("cdp-rest"),
		middleware.StripSlashes,
	)

	s.Get("/api/health", HealthInteractor())
	s.Get("/api/v1/symbol/{name}", v1.SymbolInteractor(ctx, db))
	s.Get("/api/v1/symbol/{name}/bars", v1.BarsInteractor(db))
	s.Get("/api/v1/symbols", v1.SymbolsInteractor(ctx, db))

	s.Docs("/docs", swgui.New)

	binding := fmt.Sprintf("%s:%s", cfg.HttpHost, cfg.HttpPort)
	srv := http.Server{Addr: binding, Handler: s}

	slog.Default().InfoContext(ctx, "Service running.", "binding", binding, "version", v)
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
