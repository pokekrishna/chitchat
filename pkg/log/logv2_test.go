package log

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

type mockMetaLogger struct {
	out string
}

func (mml *mockMetaLogger) basePrinter(v ...interface{}) {
	mml.out = fmt.Sprintln(v...)
}

func TestLogFunctions(t *testing.T) {
	var mockLogger *logger
	testcases := []struct {
		setLevel int
		logFunc  func(l *logger, v ...interface{})
		inputMsg string
		want     string
	}{
		{
			// TODO: use CONST in the want fields
			setLevel: 3,
			logFunc:  (*logger).Info,
			inputMsg: "foo",
			want:     "INFO: [foo]\n",
		},
		{
			setLevel: 2,
			logFunc:  (*logger).Info,
			inputMsg: "bar",
			want:     "",
		},
		{
			setLevel: 1,
			logFunc:  (*logger).Warn,
			inputMsg: "bar",
			want:     "",
		},
		{
			setLevel: 2,
			logFunc:  (*logger).Warn,
			inputMsg: "bar",
			want:     "WARN: [bar]\n",
		},
	}

	for _, tc := range testcases{
		t.Run(fmt.Sprintf("Running %v with param %s",
			runtime.FuncForPC(reflect.ValueOf(tc.logFunc).Pointer()).Name(), tc.inputMsg), func(t *testing.T) {

			mockLogger = nil

			// since the function Initialize() is not being tested here,
			// the creation of the object is customized for mocking
			mockLogger = &logger{tc.setLevel, &mockMetaLogger{}}

			tc.logFunc(mockLogger, tc.inputMsg)

			if out := mockLogger.metaLogger.(*mockMetaLogger).out; out != tc.want {
				t.Errorf("log func will print '%s', expected '%s'", out, tc.want)
			}
		})
	}
}

func TestIsInitialized(t *testing.T) {
	t.Run("When logger is not initialized, isInitialized should return false ", func(t *testing.T) {
			oldDefaultLogger := defaultLogger
			ResetForTests()
			defer func() {defaultLogger = oldDefaultLogger}()

			got := defaultLogger.isInitialized()
			expected := false
			if got != expected {
				t.Errorf("isInitialized returned %v. wanted: %v", got, expected)
			}
	})
}

func TestWithoutInitialization(t *testing.T) {
	t.Run("Logging with Error(), without Initialization, should log a static warning msg and not the input msg",
		func(t *testing.T) {
			oldDefaultLogger := defaultLogger
			ResetForTests()
			defer func() {defaultLogger = oldDefaultLogger}()

			var b bytes.Buffer
			log.SetOutput(&b)
			msg := "Dummy Error"
			Error(msg)
			got := b.String()
			expected := "Logger not Initialized. Use the Initialize function"
			if !strings.Contains(got, expected) {
				t.Errorf("got and want differ!\ngot: %s\nwant: %s", got, expected)
			}
			if strings.Contains(got, msg){
				t.Errorf("got should not contain the msg\ngot: %s\nmsg: %s", got, msg)
			}
	})
}
