package liblog

import (
	"context"
	"net/http"

	"github.com/rs/xid"
)

func Middleware(svcName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := newID()
			ctx := context.WithValue(r.Context(), cRequestId, id)
			ctx = context.WithValue(ctx, cServiceName, svcName)
			w.Header().Set(cRequestId, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HydrateServiceName(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, cServiceName, name)
}

func newID() string {
	return xid.New().String()
}
