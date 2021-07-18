package middleware

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pokekrishna/chitchat/pkg/content"
	"net/http"
)

func CheckRequestHeadersMiddleware(ctx context.Context) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := CheckAcceptHeader(ctx, r)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TODO: Complete implementation and add docs
func AddResponseHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: this setting of content type should be based on request scoped context
		w.Header().Set("Content-Type", fmt.Sprintf("%s", r.Context().Value(content.KeyAcceptContentType)))
	})
}


// TODO: complete implementation
func CheckAcceptHeader(parentCtx context.Context, r *http.Request) context.Context {
	headerVal := r.Header.Get("Accept")
	return content.ContextWithSupportedType(parentCtx, headerVal)
}