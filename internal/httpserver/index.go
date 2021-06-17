package httpserver

import (
	"fmt"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

func index(a *data.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		threads, err := a.Threads()
		if err != nil {
			log.Error("Cannot get threads", err)
		}

		s := data.NewSession(a.DB, nil)
		if ok := isValidSession(r, s); ok {
			err = generateHTML(w, threads,
				"layout.html", "private.navbar.html", "index.html")
			if err != nil {
				http.Redirect(w, r, fmt.Sprintf("/err?msg=%s", "Some problem occured"), 302)
			}
			log.Info("session validated for user email:", s.Email())
		} else {
			err = generateHTML(w, threads,
				"layout.html", "public.navbar.html", "index.html")
			if err != nil {
				http.Redirect(w, r, fmt.Sprintf("/err?msg=%s", "Some problem occured"), 302)
			}
			log.Info("invalid session")
		}
	}
}
