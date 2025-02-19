package liblog

import (
	"context"
	"log/slog"
)

func WithError(ctx context.Context, err error, msg string, args ...any) {
	slog.Default().ErrorContext(ctx, msg, append(args, "error", err)...)
}
