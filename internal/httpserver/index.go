package httpserver

import (
	"fmt"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

func index(app *data.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		threads, err := app.Threads()
		if err != nil {
			log.Error("Cannot get threads", err)
		}

		if sess, ok := isValidSession(r, app); ok {
			err = generateHTML(w, threads,
				"layout.html", "private.navbar.html", "index.html")
			if err != nil {
				http.Redirect(w, r, fmt.Sprintf("/err?msg=%s", "Some problem occured"), 302)
			}
			log.Info("session validated for user email:", sess.Email)
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
