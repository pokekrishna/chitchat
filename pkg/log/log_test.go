package log_test

import (
	"github.com/pokekrishna/chitchat/pkg/log"
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
			wantLevel: log.MaxLogLevel,
		},
		{
			name: "level 2 should initialize the package with level 2",
			setLevel:  2,
			wantLevel: 2,
		},
	}

	for _, tc := range testcases{
		t.Run(tc.name, func (t *testing.T){
			log.ResetForTests()
			log.Initialize(tc.setLevel)
			got := log.GetLevel()
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

func TestInfo(t *testing.T){
	// log.ResetForTests()
	// log.Initialize(3)

	// log.Info("foobar") should print foobar

}