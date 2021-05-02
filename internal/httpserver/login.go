package httpserver

import (
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)


func login(w http.ResponseWriter,r *http.Request) {
	t := parseTemplateFiles("login.layout.html", "public.navbar.html", "login.html")
	err := t.Execute(w, nil)
	if  err != nil {
		writeErrorToClient("Some error occurred. Please try later.", w)
		log.Error(err)
	}
}
