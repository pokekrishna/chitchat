package httpserver

import (
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

// errHandler displays msg from /err?msg=
func errHandler(w http.ResponseWriter, r *http.Request){
	// get the string from msg query parameter and display it
	queryParams := r.URL.Query()
	if _, ok := isValidSession(r); ok{
		err := generateHTML(w, queryParams.Get("msg"),
			"layout.html", "error.html", "private.navbar.html")
		if err != nil {
			log.Error("Error generating html from errorHandler", err)
		}
	} else {
		err := generateHTML(w, queryParams.Get("msg"),
			"layout.html", "error.html", "public.navbar")
		if err != nil {
			log.Error("Error generating html from errorHandler", err)
		}
	}

}
