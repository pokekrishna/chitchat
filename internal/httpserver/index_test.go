package httpserver

import (
	"bytes"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"text/template"
	"time"
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

	t.Run("GET on '/' should return HTTP 200 and the body should contain the threads", func(t *testing.T) {
		defer os.Chdir(ChangeDirForTest())
		log.Initialize(1)
		defer log.ResetForTests()

		indexTemplate := `
		<p class="lead">
			<a href="/thread/new">Start a thread</a> or join one below!
		</p>
		{{ range . }}
			<div class="panel panel-default">
				<div class="panel-heading">
					<span class="lead"> <i class="fa fa-comment-o"></i> {{ .Topic }}</span>
				</div>
				<div class="panel-body">
					Started by {{ .UserId }} - {{ .CreatedAt }} posts.
					<div class="pull-right">
						<a href="/thread/read?id={{.Uuid }}">Read more</a>
					</div>
				</div>
			</div>
		{{ end }}
		`
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
		var b bytes.Buffer
		tpl, err := template.New("expectedGeneratedHTML").Parse(indexTemplate)
		if err != nil {
			t.Errorf("Cannot Parse %s Template", tpl.Name())
		}
		if err := tpl.Execute(&b, threads); err != nil {
			t.Errorf("Cannot Execute %s Template", tpl.Name())
		}
		strippedExpectedGeneratedHtml := stripWhiteSpaces(b.String())

		if len(strippedExpectedGeneratedHtml) < len(stripWhiteSpaces(indexTemplate)) {
			t.Errorf("Some error with template. Generated length is smaller.")
			t.Fail()
		}

		rows := sqlmock.NewRows([]string{"id", "uuid", "topic", "user_id", "created_at"}).
			AddRow(threads[0].Id, threads[0].Uuid, threads[0].Topic, threads[0].UserId, threads[0].CreatedAt).
			AddRow(threads[1].Id, threads[1].Uuid, threads[1].Topic, threads[1].UserId, threads[1].CreatedAt)
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

		responseBody := w.Body.Bytes()
		strippedResponseBody := stripWhiteSpaces(string(responseBody))

		if !strings.Contains(strippedResponseBody, strippedExpectedGeneratedHtml) {
			t.Errorf("Respose does not contain the expected data.\n"+
				"Expected snippet: %s\n"+
				"Response received: %s", strippedExpectedGeneratedHtml, strippedResponseBody)
		}
	})
}

func stripWhiteSpaces(input string) string {
	return strings.Join(strings.Fields(input), "")
}
