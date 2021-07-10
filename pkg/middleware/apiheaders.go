package middleware

import "net/http"

// TODO: if doing both the things validation and addition is cumbersome in ...
// TODO: ... one middleware, split it.

// ValidateRequestHeadersAddResponseHeaders validates presence of request headers
// (keys and desired values) in the incoming request and adds common response
// headers for api
func ValidateRequestHeadersAddResponseHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// trying MVP, addition of response headers
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)

	})

}