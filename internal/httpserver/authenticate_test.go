package httpserver

import (
	"bytes"
	"github.com/pokekrishna/chitchat/internal/data"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestAuthenticate(t *testing.T){
	app, mock := NewMock()
	users := []data.User{
		data.User{
			Id:        1,
			Uuid:      "sample-uuid-1",
			Name:      "Peter",
			Email:     "p@f.com",
			Password:  data.Encrypt("pass1"),
			CreatedAt: time.Now(),
		},
		data.User{
			Id:        2,
			Uuid:      "sample-uuid-9",
			Name:      "Juan",
			Email:     "juan@foo.com",
			Password:  data.Encrypt("passjuan"),
			CreatedAt: time.Now(),
		},
	}
	userRows := mock.NewRows([]string{"id", "uuid", "name", "email", "password", "created_at"})
	for _, u := range users {
		userRows.AddRow(u.Id, u.Uuid, u.Name, u.Email, u.Password, u.CreatedAt)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/authenticate", authenticate(app))

	t.Run("unregistered user should not read password and not log in", func(t *testing.T) {

		data := url.Values{}
		data.Set("email", "random@random.com")

		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPost, "/authenticate", strings.NewReader(data.Encode()))
		if err != nil {
			t.Error("Cannot create a request", err)
		}
		mux.ServeHTTP(w, r)

		expectedStatus := http.StatusForbidden
		expectedBody := []byte("Cannot find user")
		if w.Code != expectedStatus {
			t.Errorf("Got response code %d\n Expected Response code %d", w.Code, expectedStatus)
		}

		if bytes.Compare(w.Body.Bytes(), expectedBody) != 0{
			t.Errorf("Got response body %s\n Expected response body %d", w.Body, expectedBody)
		}
		// we make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})
}
