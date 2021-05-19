package log

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func TestWithoutInitialization(t *testing.T) {
	// log.ResetForTests()

	// log.Info("foobar") should not write anything on the stream/buffer
}

func TestInitialize(t *testing.T) {
	testcases := []struct{
		name string
		setLevel int
		wantLevel int
	}{
		{
			name : "negative level should keep the package uninitialized",
			setLevel:  -1,
			wantLevel: 0,
		},
		{
			name: "level 0 should keep the package uninitialized",
			setLevel:  0,
			wantLevel: 0,
		},
		{
			name: "level beyond max level should initialize the package with max level",
			setLevel:  999,
			wantLevel: MaxLogLevel,
		},
		{
			name: "level 2 should initialize the package with level 2",
			setLevel:  2,
			wantLevel: 2,
		},
	}

	for _, tc := range testcases{
		t.Run(tc.name, func (t *testing.T){
			ResetForTests()
			Initialize(tc.setLevel)
			got := GetLevel()
			if tc.wantLevel != got {
				t.Errorf("Initalize(%d) resulted in log level '%d'. Expected level '%d'",
					tc.setLevel, got, tc.wantLevel)
			}
		})
	}
}

func TestInitializeSingleton(t *testing.T){
	// more than once should print warning and level should not change

}

func TestLogFunctions(t *testing.T){
	testcases := []struct{
		setLevel int

		// TODO: Cons of using this way of mocking. Have to use a generic func definition,
		// instead of a type
		logFunc func(...interface{})

		inputMsg string
		want string
	}{
		{
			// TODO: use CONST in the want fields
			setLevel: 3,
			logFunc: Info,
			inputMsg: "foo",
			want: "INFO: [foo]\n",
		},
		{
			setLevel: 2,
			logFunc: Info,
			inputMsg: "bar",
			want: "",
		},
		{
			setLevel: 1,
			logFunc: Warn,
			inputMsg: "bar",
			want: "",
		},
		{
			setLevel: 2,
			logFunc: Warn,
			inputMsg: "bar",
			want: "WARN: [bar]\n",
		},
	}

	var output string
	var mockPrinter BasePrinter = func(v ...interface{}){
		t.Logf("mock Printer called with (%v)", v)
		output = fmt.Sprintln(v...)
	}
	for _, tc := range testcases{
		t.Run(fmt.Sprintf("Running %v with param %s",
			runtime.FuncForPC(reflect.ValueOf(tc.logFunc).Pointer()).Name(), tc.inputMsg), func(t *testing.T) {
			output = ""
			ResetForTests()
			Initialize(tc.setLevel)

			// TODO: cons of this way of testing.
			// you cannot mock from outside this packaage.
			// relying on private var.
			printer = mockPrinter
			tc.logFunc(tc.inputMsg)

			if output != tc.want{
				t.Errorf("log func will print '%s', expected '%s'", output, tc.want)
			}
		})
	}
}