package data_test

import (
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/enhancederror"
	"testing"
)

var userTest = []struct {
	u              data.user
	expected error
}{
	{data.user{
		Name:     "Peter Jones",
		Email:    "peter@gmail.com",
		Password: "peter_pass",
	},
	nil,
	},
	{
		data.user{
			Name:     "John Smith",
			Email:    "john@gmail.com",
			Password: "",
		},
		&data.InvalidUser{Reason: "Password not set"},
	},
}

func TestUser_Validate(t *testing.T) {
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
	for _, ut := range userTest{
		if err := ut.u.Create(); !enhancederror.IsEqual(err, ut.expected) {
			t.Errorf("%v.Create() returned %v, expected %v.",
				ut.u, err, ut.expected)
		}
	}
}

