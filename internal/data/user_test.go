package data_test

import (
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/enhancederror"
	"testing"
	"time"
)

func TestUser_Validate(t *testing.T) {
	userTest := []struct {
		u           *data.User
		expectedErr error
	}{
		{
			&data.User{
				Id:        1,
				Uuid:      "sample-uuid-1",
				Name:      "Peter Jones",
				Email:     "peter@gmail.com",
				Password:  "peter_pass",
				CreatedAt: time.Now(),
			},
			nil,
		},
		{
			&data.User{
				Name:     "John Smith",
				Email:    "john@gmail.com",
				Password: "",
			},
			&data.InvalidUser{Reason: "Password not set"},
		},
	}

	for _, ut := range userTest {
		got := ut.u.Validate()
		if !enhancederror.IsEqual(got, ut.expectedErr) {
			t.Errorf("User validation failed for user: %v", ut.u)
			t.Errorf("%v.Validate() returned %v, expected %v.",
				ut.u, got, ut.expectedErr)
		}
	}
}

func TestApp_CreateUser(t *testing.T) {
	app, mock := NewMock()
	userTest := []struct {
		u           *data.User
		expectedErr error
	}{
		{
			&data.User{
				Id:        1,
				Uuid:      "sample-uuid-1",
				Name:      "Peter Jones",
				Email:     "peter@gmail.com",
				Password:  "peter_pass",
				CreatedAt: time.Now(),
			},
			nil,
		},
		{
			&data.User{
				Name:     "John Smith",
				Email:    "john@gmail.com",
				Password: "",
			},
			&data.InvalidUser{Reason: "Password not set"},
		},
	}

	for _, ut := range userTest {
		if ut.expectedErr == nil {
			//creating expected rows
			rows := mock.NewRows([]string{"id", "uuid", "name", "email", "password", "created_at"}).
				AddRow(userTest[0].u.Id, userTest[0].u.Uuid, userTest[0].u.Name, userTest[0].u.Email, userTest[0].u.Password,
					userTest[0].u.CreatedAt)

			ePrep := mock.ExpectPrepare("^INSERT INTO users (.+) VALUES (.+) RETURNING .+$")
			eQuery := ePrep.ExpectQuery()
			eQuery.WillReturnRows(rows)
		}
		if err := app.CreateUser(ut.u); !enhancederror.IsEqual(err, ut.expectedErr) {
			t.Errorf("Expectations Mismatch from CreateUser(%v)\n Got: %v\n Expected: %v.",
				ut.u, err, ut.expectedErr)
		}
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestApp_CreateSession(t *testing.T) {
	app, mock := NewMock()
	sessionTest := []struct {
		desc        string
		s           *data.Session
		expectedErr error
	}{
		{
			"InValid Session - incomplete user should return error",
			&data.Session{
				Id:        1,
				Uuid:      "sample-uuid-1",
				Email:     "peter@gmail.com",
				User: &data.User{},
				CreatedAt: time.Now(),
			},
			&data.InvalidUser{Reason: "Empty Name"},
		},
		{
			"Valid Session - should return nil error",
			&data.Session{
				Id:        1,
				Uuid:      "sample-uuid-1",
				Email:     "",
				User: &data.User{
					Id:        10,
					Name:      "Peter",
					Email:     "peter@gmail.com",
					Password:  "peter_pass",
				},
				CreatedAt: time.Now(),
			},
			nil,
		},
	}

	for _, st := range sessionTest {
		t.Run(st.desc, func(t *testing.T) {
			if st.expectedErr == nil {
				//creating expected rows
				rows := mock.NewRows([]string{"id", "uuid", "email", "created_at"}).
					AddRow(sessionTest[0].s.Id, sessionTest[0].s.Uuid, sessionTest[0].s.User.Email,
						 sessionTest[0].s.CreatedAt)

				ePrep := mock.ExpectPrepare("^INSERT INTO sessions (.+) VALUES (.+) RETURNING .+$")
				eQuery := ePrep.ExpectQuery()
				eQuery.WillReturnRows(rows)
			}
			if err := app.CreateSession(st.s); !enhancederror.IsEqual(err, st.expectedErr) {
				t.Errorf("Expectations Mismatch from CreateSession(%v)\n Got: %v\n Expected: %v.",
					st.s, err, st.expectedErr)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
