package data

import "testing"

func TestInitialize(t *testing.T){
	if db, err := Initialize(); err != nil {
		t.Error("Cannot Initialize connection to db.", err)
	}else {
		if db == nil{
			t.Errorf("db should NOT be nil after calling Initialize()")
		}
	}
}

