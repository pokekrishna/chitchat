package httpserver

import (
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

var cookieName = "_cookie"

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
			Name: cookieName,
			Value: s.Uuid,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/", 302)
	} else {
		http.Redirect (w, r, "/login", 302)
	}
}

// Check if the session id is valid
func isValidSession(r *http.Request) (session *data.Session, ok bool){
	ok = false
	// get the session from the cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Info("Cannot find cookie from the request object")
		return
	}

	// find the session object by the uuid from the cookie
	session, err = data.SessionByUuid(cookie.Value)
	if err != nil {
		// there was a problem in finding the session, hence invalid session
		// Log the error and return accordingly
		log.Warn("Cannot find session by Uuid:", cookie.Value)
		return
	} else {
		ok = true
	}
	return
}