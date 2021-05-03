package httpserver

import (
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

func authenticate(w http.ResponseWriter, r *http.Request) {
	// user authentication here with DB reconciliation
	if err := r.ParseForm(); err != nil {
		log.Error("Cannot parse form", err)
	}

	u, err := data.UserByEmail(r.PostFormValue("email"))
	if err != nil {
		msg := "Cannot find user"
		log.Info(msg, err)
		writeErrorToClient(msg, w)
	}

	if u.Password == data.Encrypt(r.PostFormValue("password")) {
		log.Info("Authentication successful for user email", u.Email)
		// create a session
		s, err := u.CreateSession()
		if err != nil {
			log.Error("Cannot create session", err)
		}

		// create a cookie based on session
		cookie := http.Cookie{
			Name: "_cookie",
			Value: s.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect (w, r, "/login", 302)
	}
}