package liblog

import (
	"context"
	"net/http"

	"github.com/rs/xid"
)

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := newID()
			ctx := context.WithValue(r.Context(), cRequestId, id)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func HydrateServiceName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, cServiceName, name)
}

func newID() string {
	return xid.New().String()
}
