package data_test

import (
	"github.com/pokekrishna/chitchat/internal/data"
	"testing"
)

func TestMain(m *testing.M){
	err := data.Initialize()
	if err != nil {
		panic("cannot initialize db")
	}

	cleanDB()
	m.Run()
}

func TestCreateUUID(t *testing.T) {
	if data.CreateUUID() == "" {
		t.Error("CreateUUID() returned empty string.")
	}
}


func TestEncrypt(t *testing.T) {
	got := data.Encrypt("hello")
	expected := "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"
	if got != expected {
		t.Errorf("got %s, expected %s", got, expected)
	}
}

func cleanDB(){
	deleteAllSessions()
}

func deleteAllSessions() {
}