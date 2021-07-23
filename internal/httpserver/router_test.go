package httpserver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TODO: interim func name, resolve name conflict by better design
func NewTestData() (*sql.DB, sqlmock.Sqlmock, []*data.Thread, []byte, []byte) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("error instantiating sqlmock %s", err))
	}

	timeT1 := time.Now()
	timeT2 := time.Now()
	threads := []*data.Thread{
		&data.Thread{
			Id:        1,
			Uuid:      "uuid-sample-1",
			Topic:     "topic1",
			UserId:    2,
			CreatedAt: timeT1,
		},
		&data.Thread{
			Id:        2,
			Uuid:      "uuid-sample-2",
			Topic:     "topic2",
			UserId:    5,
			CreatedAt: timeT2,
		},
	}

	timeT1RFC3339, err := timeT1.MarshalJSON()
	if err != nil {
		panic(fmt.Sprint("cannot convert time stamp to RFC3339 time", err))
	}

	timeT2RFC3339, err := timeT2.MarshalJSON()
	if err != nil {
		panic(fmt.Sprint("cannot convert time stamp to RFC3339 time", err))
	}

	return db, mock, threads, timeT1RFC3339, timeT2RFC3339
}

func TestGETThreads (t *testing.T){
	log.Initialize(3)
	defer log.ResetForTests()

	t.Run("GET on '/threads' should return HTTP 200 and json array response", func(t *testing.T) {
		db, mock, threads, timeT1RFC3339, timeT2RFC3339  := NewTestData()

		jsonThreads := fmt.Sprintf("[{\"Id\":1,\"Uuid\":\"uuid-sample-1\",\"Topic\":\"topic1\",\"UserId\":2,\"CreatedAt\":%s},{\"Id\":2,\"Uuid\":\"uuid-sample-2\",\"Topic\":\"topic2\",\"UserId\":5,\"CreatedAt\":%s}]",
			timeT1RFC3339, timeT2RFC3339)

		rows := sqlmock.NewRows([]string{"id", "uuid", "topic", "user_id", "created_at"}).
			AddRow(threads[0].Id, threads[0].Uuid, threads[0].Topic, threads[0].UserId, threads[0].CreatedAt).
			AddRow(threads[1].Id, threads[1].Uuid, threads[1].Topic, threads[1].UserId, threads[1].CreatedAt)
		mock.ExpectQuery("^SELECT (.+) FROM threads order by created_at desc$").WillReturnRows(rows)

		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, "/api/v1/threads", nil)
		if err != nil {
			t.Error("Cannot create a request", err)
		}
		mux := Router(context.Background(), db)
		mux.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Code, "Response code should be equal.")
		assert.Equal(t, jsonThreads, string(w.Body.Bytes()), "Response Body should be equal")

		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})
}
