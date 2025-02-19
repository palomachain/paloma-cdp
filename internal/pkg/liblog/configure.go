package liblog

import (
	"context"
	"log/slog"
	"os"
)

const (
	cServiceName = "x-cdp-service-name"
	cRequestId   = "x-request-id"
)

var logFields = []string{cServiceName, cRequestId}

type contextHandler struct {
	slog.Handler
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, key := range logFields {
		if value, ok := ctx.Value(key).(string); ok {
			r.Add(key, value)
		}
	}

	return h.Handler.Handle(ctx, r)
}

func Configure() {
	baseHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(&contextHandler{Handler: baseHandler})
	slog.SetDefault(logger)
}
