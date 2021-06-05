package data_test

import (
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/enhancederror"
	"testing"
)

func TestUser_Validate(t *testing.T) {
	u1 := data.NewUser(db)
	u1.SetName("Peter Jones")
	u1.SetEmail("peter@gmail.com")
	u1.SetPassword("peter_pass")

	u2 := data.NewUser(db)
	u2.SetName("John Smith")
	u2.SetEmail("john@gmail.com")
	u2.SetPassword("")
	var userTest = []struct {
		u              data.UserInterface
		expected error
	}{
		{u1,
			nil,
		},
		{
			u2,
			&data.InvalidUser{Reason: "Password not set"},
		},
	}

	for _,ut := range userTest{
		got := ut.u.Validate()
		if !enhancederror.IsEqual(got, ut.expected){
			t.Errorf("User validation failed for user: %v", ut.u)
			t.Errorf("%v.Validate() returned %v, expected %v.",
				ut.u, got, ut.expected)
		}
	}
}

func TestUser_Create(t *testing.T) {
	u1 := data.NewUser(db)
	u1.SetName("Peter Jones")
	u1.SetEmail("peter@gmail.com")
	u1.SetPassword("peter_pass")

	u2 := data.NewUser(db)
	u2.SetName("John Smith")
	u2.SetEmail("john@gmail.com")
	u2.SetPassword("")
	var userTest = []struct {
		u              data.UserInterface
		expected error
	}{
		{u1,
			nil,
		},
		{
			u2,
			&data.InvalidUser{Reason: "Password not set"},
		},
	}

	for _, ut := range userTest{
		if err := ut.u.Create(); !enhancederror.IsEqual(err, ut.expected) {
			t.Errorf("%v.Create() returned %v, expected %v.",
				ut.u, err, ut.expected)
		}
		// TODO: select from DB to check the inserted value
	}
}

