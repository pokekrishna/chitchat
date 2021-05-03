package httpserver

import (
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request){
	threads, err := data.Threads()
	if err != nil {
		log.Error("Cannot get threads", err)
	}

	if isValidSession(r){
		err = generateHTML(w, threads,
			"layout.html", "private.navbar.html", "index.html")
		if err != nil {
			// TODO: respond using err handler
		}
	} else {
		err = generateHTML(w, threads,
			"layout.html", "public.navbar.html", "index.html")
		if err != nil {
			// TODO: respond using err handler
		}

	}

}
