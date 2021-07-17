package middleware

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pokekrishna/chitchat/pkg/content"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckRequestHeadersMiddleware(t *testing.T) {
	t.Run("Request with header 'Accept: foo' should have content.TypeNotSupported enum passed on in the req context",
		func(t *testing.T) {
			headerVal := "foo"
			router := mux.NewRouter()
			router.Use(CheckRequestHeadersMiddleware(context.Background()))
			var captureRequest *http.Request
			router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
				captureRequest = request
			})

			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Error("Cannot create request", err)
			}
			r.Header.Set("Accept", headerVal)
			router.ServeHTTP(w, r)
			assert.Equal(t, content.TypeNotSupported, captureRequest.Context().Value(content.KeyAcceptContentType))
		})
	t.Run("Request with no 'Accept' header should have content.TypeNotSupported enum passed on in the req context",
		func(t *testing.T) {
			router := mux.NewRouter()
			router.Use(CheckRequestHeadersMiddleware(context.Background()))
			var captureRequest *http.Request
			router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
				captureRequest = request
			})

			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Error("Cannot create request", err)
			}
			router.ServeHTTP(w, r)
			assert.Equal(t, content.TypeNotSupported, captureRequest.Context().Value(content.KeyAcceptContentType))
		})
	t.Run("Request with header 'Accept: json' should have content.TypeJSON enum passed on in the req context",
		func(t *testing.T) {
			headerVal := "json"
			router := mux.NewRouter()
			router.Use(CheckRequestHeadersMiddleware(context.Background()))
			var captureRequest *http.Request
			router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
				captureRequest = request
			})

			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Error("Cannot create request", err)
			}
			r.Header.Set("Accept", headerVal)
			router.ServeHTTP(w, r)
			assert.Equal(t, content.TypeJSON, captureRequest.Context().Value(content.KeyAcceptContentType))
		})
	t.Run("Request with header 'Accept: jSOn' (caseinsensitivity) should have content.TypeJSON enum passed on in the req context",
		func(t *testing.T) {
			headerVal := "jSOn"
			router := mux.NewRouter()
			router.Use(CheckRequestHeadersMiddleware(context.Background()))
			var captureRequest *http.Request
			router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
				captureRequest = request
			})

			w := httptest.NewRecorder()
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Error("Cannot create request", err)
			}
			r.Header.Set("Accept", headerVal)
			router.ServeHTTP(w, r)
			assert.Equal(t, content.TypeJSON, captureRequest.Context().Value(content.KeyAcceptContentType))
		})
}

