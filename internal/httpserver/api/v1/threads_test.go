package v1

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TODO: NewMock is code duplication in multiple packages
// NewMock instantiates mock elements necessary for testing.
func NewMock() (*data.App, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("error instantiating sqlmock %s", err))
	}
	app := &data.App{DB: db}
	return app, mock
}

func TestThreads(t *testing.T){
	app, mock := NewMock()
	mux := http.NewServeMux()
	mux.HandleFunc("/", Threads(app))

	t.Run("GET on '/threads' should return HTTP 200", func(t *testing.T) {
		log.Initialize(1)
		defer log.ResetForTests()

		threads := []data.Thread{
			data.Thread{
				Id:        1,
				Uuid:      "uuid-sample-1",
				Topic:     "topic1",
				UserId:    2,
				CreatedAt: time.Now(),
			},
			data.Thread{
				Id:        2,
				Uuid:      "uuid-sample-2",
				Topic:     "topic2",
				UserId:    5,
				CreatedAt: time.Now(),
			},
		}

		rows := sqlmock.NewRows([]string{"id", "uuid", "topic", "user_id", "created_at"}).
			AddRow(threads[0].Id, threads[0].Uuid, threads[0].Topic, threads[0].UserId, threads[0].CreatedAt).
			AddRow(threads[1].Id, threads[1].Uuid, threads[1].Topic, threads[1].UserId, threads[1].CreatedAt)
		mock.ExpectQuery("^SELECT (.+) FROM threads order by created_at desc$").WillReturnRows(rows)

		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/threads", nil)
		if err != nil {
			t.Error("Cannot create a request", err)
		}

		// mux.ServeHTTP helps sending the request to the handler without
		// running a HTTP server.
		mux.ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("Response code not", http.StatusOK)
		}
		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})
}
