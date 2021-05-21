package log

import (
	"fmt"
	"testing"
)

type mockMetaLogger struct {
	out string
}

func (mml *mockMetaLogger) basePrinter(v ...interface{}){
	mml.out = fmt.Sprintln(v...)
}

func TestPlayground(t *testing.T){
	mockLogger := &logger{3, &mockMetaLogger{}}
	mockLogger.Info("foo")

	if out := mockLogger.metaLogger.(*mockMetaLogger).out; out != "INFO: [foo]\n"{
		t.Errorf("wanted 'INFO: [foo]\n', got '%s'", out)
	}
}