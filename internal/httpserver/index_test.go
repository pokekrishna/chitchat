package httpserver

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pokekrishna/chitchat/internal/data"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestIndex(t *testing.T) {
	//m := &mockThread{}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error opening a stud db connection %s", err)
	}
	app := &data.App{DB: db}
	mux := http.NewServeMux()
	mux.HandleFunc("/", index(app))

	t.Run("GET on '/' should return HTTP 200", func(t *testing.T) {
		// TODO : move this DIR related acts to a function
		cwd, err := os.Getwd()
		if err != nil {
			t.Errorf("Cannot get current working directory")
			t.Fail()
		}

		if err := os.Chdir("../.."); err != nil {
			t.Errorf("Cannot change directory")
			t.Fail()
		}
		defer os.Chdir(cwd)

		rows := sqlmock.NewRows([] string {"id", "uuid", "topic", "user_id", "created_at"}).
			AddRow(1, "uuid-sample-1", "topic1", 2, "TIME1").
			AddRow(2, "uuid-sample-2", "topic2", 5, "TIME2")
		mock.ExpectQuery("^SELECT (.+) FROM threads order by created_at desc$").WillReturnRows(rows)

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

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}