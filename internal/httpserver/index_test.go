package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	t.Run("GET on '/' should return HTTP 200", func(t *testing.T) {
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Error("Cannot create a request", err)
		}

		// mux.ServeHTTP helps sending the request to the handler without
		// running a HTTP server.
		mux.ServeHTTP(w, r)

		if w.Code != http.StatusOK {
			t.Error("Response code not", http.StatusOK)
		}
	})

}

