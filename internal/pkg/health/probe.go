package health

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
)

func StartHealthProbe(ctx context.Context, addr, endpoint string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			liblog.WithError(ctx, err, "Failed to shutdown health probe server.")
		}
	}()

	go func() {
		slog.Default().InfoContext(ctx, "Starting health probe.", "healthz-address", addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			liblog.WithError(ctx, err, "Health probe listen error.")
		}
	}()
}
