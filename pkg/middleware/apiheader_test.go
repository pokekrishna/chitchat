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
	t.Run("Request with header 'Accept: foo' should have content.TypeNotSupported passed on in the req context",
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
			captureRequestContext := captureRequest.Context().(*content.Context)
			assert.Equal(t, content.TypeNotSupported, captureRequestContext.ContentType())
		})
	t.Run("Request with no 'Accept' header should have content.TypeNotSupported passed on in the req context",
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
			captureRequestContext := captureRequest.Context().(*content.Context)
			assert.Equal(t, content.TypeNotSupported, captureRequestContext.ContentType())
		})
	t.Run("Request with header 'Accept: application/json' should have content.TypeJSON passed on in the req context",
		func(t *testing.T) {
			headerVal := "application/json"
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
			captureRequestContext := captureRequest.Context().(*content.Context)
			assert.Equal(t, content.TypeJSON, captureRequestContext.ContentType())
		})
	t.Run("Request with header 'Accept: application/jSOn' (caseinsensitivity) should have content.TypeJSON passed on in the req context",
		func(t *testing.T) {
			headerVal := "application/jSOn"
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
			captureRequestContext := captureRequest.Context().(*content.Context)
			assert.Equal(t, content.TypeJSON, captureRequestContext.ContentType())
		})
}

func TestCheckAcceptHeader(t *testing.T) {
	testcases := []struct{
		description         string
		headerKey           string
		headerVals          []string
		expectedContentType string
	}{
		{
			headerKey:           "",
			headerVals:          []string{""},
			expectedContentType: content.TypeNotSupported,
			description:         "No Accept header should not be supported but Context Key should be set",
		},
		{
			headerKey:           "Accept",
			headerVals:          []string{"foo"},
			expectedContentType: content.TypeNotSupported,
			description:         "Accept header foo should not be supported",
		},
		{
			headerKey:           "accept",
			headerVals:          []string{"json"},
			expectedContentType: content.TypeNotSupported,
			description:         "Accept header with non standard value is not supported",
		},
		{
			headerKey:           "accept",
			headerVals:          []string{"application/json"},
			expectedContentType: content.TypeJSON,
			description:         "Accept header with standard value is supported",
		},
		{
			headerKey:           "Accept",
			headerVals:          []string{"application/json", "xml"},
			expectedContentType: content.TypeJSON,
			description:         "Accept header with one supported and one non supported type should consider the first",
		},
		{
			headerKey:           "accept",
			headerVals:          []string{"text", "application/json", "xml"},
			expectedContentType: content.TypeNotSupported,
			description:         "Accept header with one supported and one non supported type should consider the first",
		},

	}

	for _, tc := range testcases{
		t.Run(tc.description, func(t *testing.T) {
			r, err := http.NewRequest("", "", nil)
			if err != nil {
				t.Error("Cannot create request", err)
			}
			for _, hv := range tc.headerVals {
				r.Header.Add(tc.headerKey, hv)
			}
			ctx := CheckAcceptHeader(context.Background(), r).(*content.Context)
			assert.Equal(t, tc.expectedContentType, ctx.ContentType())
		})
	}

}

func TestAddResponseHeadersMiddleware(t *testing.T) {
	testcases := []struct{
		description             string
		headerVal              	string
		responseBody			string
		expectedRespContentType string
		expectedRespCode int
	}{
		{
			description:            "Unsupported 'Accept' request header should result in Content-Type=text/plain response header",
			headerVal:              "",
			responseBody: 			"foobar",
			expectedRespContentType: "text/plain; charset=utf-8",
			expectedRespCode: http.StatusOK,
		},
	}
	for _, tc := range testcases{
		t.Run(tc.description, func(t *testing.T) {
			router := mux.NewRouter()
			router.Use(AddResponseHeadersMiddleware)
			router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(tc.responseBody))
			})

			w := httptest.NewRecorder()
			rCtx := content.ContextWithSupportedType(context.Background(), tc.headerVal)
			r, err := http.NewRequestWithContext(rCtx, http.MethodGet, "/", nil)
			if err != nil {
				t.Error("Cannot create request", err)
			}
			router.ServeHTTP(w, r)
			assert.Equal(t, tc.expectedRespContentType, w.Header().Get("Content-type"))
			assert.Equal(t, tc.responseBody, w.Body.String())
			assert.Equal(t, tc.expectedRespCode, w.Code)
		})
	}

	t.Run("content.Context not set in the request context, should return code 415 with no content type ",
		func(t *testing.T) {
			responseBody := "foobar"
			expectedRespContentType := ""
			expectedRespCode := http.StatusUnsupportedMediaType

			router := mux.NewRouter()
			router.Use(AddResponseHeadersMiddleware)
			router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(responseBody))
			})

			w := httptest.NewRecorder()
			r, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
			if err != nil {
				t.Error("Cannot create request", err)
			}
				router.ServeHTTP(w, r)
				assert.Equal(t, expectedRespContentType, w.Header().Get("Content-type"))
				assert.Equal(t, responseBody, w.Body.String())
				assert.Equal(t, expectedRespCode, w.Code)
		})
}

