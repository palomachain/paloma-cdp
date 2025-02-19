package rest

import (
	"log"
	"net/http"

	v1 "github.com/palomachain/paloma-cdp/internal/app/rest/v1"
	"github.com/swaggest/openapi-go/openapi31"
	"github.com/swaggest/rest/response/gzip"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v5emb"
)

func Run() {
	s := web.NewService(openapi31.NewReflector())

	// Init API documentation schema.
	s.OpenAPISchema().SetTitle("Basic Example")
	s.OpenAPISchema().SetDescription("This app showcases a trivial REST API.")
	s.OpenAPISchema().SetVersion("v1.2.3")

	// Setup middlewares.
	s.Wrap(
		gzip.Middleware, // Response compression with support for direct gzip pass through.
	)

	// Add use case handler to router.
	s.Get("/api/v1/symbol/{name}", v1.SymbolInteractor())
	s.Get("/api/v1/symbol/{name}/bars", v1.BarsInteractor())
	s.Get("/api/v1/symbols", v1.SymbolsInteractor())

	// Swagger UI endpoint at /docs.
	s.Docs("/docs", swgui.New)

	// Start server.
	log.Println("http://localhost:8011/docs")
	if err := http.ListenAndServe("localhost:8011", s); err != nil {
		log.Fatal(err)
	}
}
