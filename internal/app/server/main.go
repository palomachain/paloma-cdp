package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/palomachain/paloma-cdp/internal/app/config"
	"github.com/palomachain/paloma-cdp/internal/app/gql"
	"github.com/palomachain/paloma-cdp/internal/app/gql/resolvers"
	"github.com/palomachain/paloma-cdp/internal/app/rest"
	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/vektah/gqlparser/v2/ast"
)

func Run() {
	os.Setenv("CDP_PSQL_ADDRESS", "localhost:5432")
	os.Setenv("CDP_PSQL_USER", "cdp")
	os.Setenv("CDP_PSQL_PASSWORD", "trustno1")
	os.Setenv("CDP_PSQL_DATABASE", "cdp")
	os.Setenv("CDP_GQL_PORT", "8080")

	rest.Run()
	ctx, _ := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, "request_id", "hurensohn")
	liblog.Configure()

	cfg, err := config.Parse()
	if err != nil {
		slog.Default().ErrorContext(ctx, "failed to parse config: %v", err)
		panic(err)
	}

	db, err := persistence.New(ctx, &cfg.Persistence)
	if err != nil {
		slog.Default().ErrorContext(ctx, "failed to connect to database: %v", err)
		panic(err)
	}

	if err := db.Migrate(ctx); err != nil {
		slog.Default().ErrorContext(ctx, "failed to migrate database: %v", err)
		panic(err)
	}

	router := chi.NewRouter()
	router.Use(liblog.Middleware())

	srv := createServer(db)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	slog.Default().InfoContext(ctx, "connect to http://localhost:%s/ for GraphQL playground", cfg.GraphQL.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.GraphQL.Port, router))
}

func createServer(db *persistence.Database) *handler.Server {
	srv := handler.New(gql.NewExecutableSchema(gql.Config{
		Resolvers: &resolvers.Resolver{
			Db: db,
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}
