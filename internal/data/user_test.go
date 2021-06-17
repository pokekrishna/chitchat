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
				Id:       1,
				Uuid:	"sample-uuid-1",
				Name:     "Peter Jones",
				Email:    "peter@gmail.com",
				Password: "peter_pass",
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

	for _,ut := range userTest{
		got := ut.u.Validate()
		if !enhancederror.IsEqual(got, ut.expectedErr){
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
				Id:       1,
				Uuid:	"sample-uuid-1",
				Name:     "Peter Jones",
				Email:    "peter@gmail.com",
				Password: "peter_pass",
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

	for _, ut := range userTest{
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