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

// AddResponseHeadersMiddleware gets the requested content type from the request
// context, set previously by CheckRequestHeadersMiddleware, and sets it in the
// response header. In any case if request case request content type did not map
// to support content type, do not set response content type.
//
// next Handler can override the response content type header only in dire need.
// It is not advised as a general practice though.
func AddResponseHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respContentType, err := content.ExtractContentType(r)
		if err != nil {
			// do not set content-type but set error code
			w.WriteHeader(http.StatusUnsupportedMediaType)
		} else {
			w.Header().Set("Content-Type", respContentType)
		}
		next.ServeHTTP(w, r)
	})
}


// TODO: complete implementation
func CheckAcceptHeader(parentCtx context.Context, r *http.Request) context.Context {
	headerVal := r.Header.Get("Accept")
	return content.ContextWithSupportedType(parentCtx, headerVal)
}