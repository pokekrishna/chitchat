package data_test

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"testing"
)

func init(){
	log.Initialize(3)
}
// NewMock instantiates mock elements necessary for testing.
func NewMock() (*data.App, sqlmock.Sqlmock){
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(fmt.Sprintf("error instantiating sqlmock %s", err))
	}
	app := &data.App{DB: db}
	return app, mock
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
