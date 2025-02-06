package liblog

import (
	"context"
	"log/slog"
	"os"
)

const (
	cRequestId = "x-request-id"
)

type contextHandler struct {
	slog.Handler
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value(cRequestId).(string); ok {
		r.Add(cRequestId, requestID)
	}

	return h.Handler.Handle(ctx, r)
}

func Configure() {
	baseHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(&contextHandler{Handler: baseHandler})
	slog.SetDefault(logger)
}
