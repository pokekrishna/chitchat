package data

import "testing"

func TestInitialize(t *testing.T){
	if db != nil {
		t.Errorf("db should be nil before calling Initialize()")
	}

	if err := Initialize(); err != nil {
		t.Error("Cannot Initialize connection to DB.", err)
	}else {
		if db == nil{
			t.Errorf("db should NOT be nil after calling Initialize()")
		}
	}
}
