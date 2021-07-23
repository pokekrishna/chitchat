package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pokekrishna/chitchat/pkg/content"
	"net/http"
)

// CheckRequestHeadersMiddleware assesses the 'Accept' request Header and adds
// the suitable content type in context so that the handlers do not have assess
// the content type and can make quick decisions by using the helper methods
// provided by content.Context.
//
// This middleware that scans the incoming request has a response counterpart too
// AddResponseHeadersMiddleware which adds the response header automatically by
// assessing the (*content.context).ContentType()
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
		respContentType := r.Context().Value(content.KeyAcceptContentType)
		switch respContentType :=  respContentType.(type){
		case string:
			w.Header().Set("Content-Type", respContentType)
		default:
			// do not set content-type
		}
		next.ServeHTTP(w, r)
	})
}


// TODO: complete implementation
func CheckAcceptHeader(parentCtx context.Context, r *http.Request) context.Context {
	headerVal := r.Header.Get("Accept")
	return content.ContextWithSupportedType(parentCtx, headerVal)
}