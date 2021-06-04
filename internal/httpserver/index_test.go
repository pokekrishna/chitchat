package httpserver

import (
	"database/sql"
	"github.com/pokekrishna/chitchat/internal/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockThread struct{}
func (m *mockThread) FetchAll() ([]data.ThreadInterface, error) {
	var threads []data.ThreadInterface
	threads = append(threads, &mockThread{}, &mockThread{})
	return threads, nil
}

func (m *mockThread) GetDB() *sql.DB{
	return nil
}

func TestIndex(t *testing.T) {
	m := &mockThread{}
	mux := http.NewServeMux()
	mux.HandleFunc("/", index(m))

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