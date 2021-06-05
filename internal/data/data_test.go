package data_test

import (
	"database/sql"
	"github.com/pokekrishna/chitchat/internal/data"
	"testing"
)

// db is a a global variable declared for the data_test package.
// Although use of global variable is a design flaw, it is used
// here because the `go test` utility does not allow more than
// one parameter to the Test... functions other than the *testing.T
// itself.
var db *sql.DB

func TestMain(m *testing.M){
	// create DB connection but not test it
	var err error
	db, err = data.Initialize()
	if err != nil {
		panic("cannot initialize db")
	}

	cleanDB()
	defer cleanDB()
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
	u := data.NewUser(db)
	s := data.NewSession(db, u)
	s.DeleteAllSessions()
	u.DeleteAllUsers()
}