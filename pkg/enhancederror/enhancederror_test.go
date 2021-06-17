package enhancederror_test

import (
	"github.com/pokekrishna/chitchat/pkg/enhancederror"
	"testing"
)
import "errors"

type testcase struct {
	source error
	target error
	expect bool
}

type CustomErr struct {
	msg string
}

func (c *CustomErr) Error() string {
	return c.msg
}

var testcases = []testcase{
	{
		source: nil,
		target: nil,
		expect: true,
	},

	{
		source: errors.New("flag1"),
		target: nil,
		expect: false,
	},

	{
		source: nil,
		target: errors.New("flag2"),
		expect: false,
	},

	{
		source: errors.New("flag3"),
		target: errors.New("flag4"),
		expect: false,
	},

	{
		source: &CustomErr{msg: "flag1"},
		target: nil,
		expect: false,
	},

	{
		source: nil,
		target: &CustomErr{msg: "flag1"},
		expect: false,
	},

	{
		source: &CustomErr{msg: "flag1"},
		target: &CustomErr{msg: "flag2"},
		expect: false,
	},

	{
		source: &CustomErr{msg: "flag3"},
		target: &CustomErr{msg: "flag3"},
		expect: true,
	},
}

func TestIsEqual(t *testing.T) {
	for _, tc := range testcases {
		got := enhancederror.IsEqual(tc.source, tc.target)
		if got != tc.expect {
			t.Errorf("IsEqual(%v,%v) returned %v, expected %v.",
				tc.source, tc.target, got, tc.expect)
		}
	}
}
