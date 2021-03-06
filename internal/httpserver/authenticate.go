package httpserver

import (
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
)

const (
	cookieName = "_cookie"
)

func authenticate(app *data.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// user authentication here with db reconciliation
		if err := r.ParseForm(); err != nil {
			log.Error("Cannot parse form", err)
		}

		u := &data.User{
			Email: r.PostFormValue("email"),
		}
		err := app.FindUserByEmail(u)
		if err != nil {
			msg := "Cannot find user"
			log.Info(msg, err)
			writeErrorToClient(msg, w, http.StatusForbidden)
		}

		if u.Password == data.Encrypt(r.PostFormValue("password")) {
			log.Info("Authentication successful for user email", u.Email)
			// create a session
			sess := &data.Session{User: u}
			err := app.CreateSession(sess)
			if err != nil {
				log.Error("Cannot create session", err)
			}

			// create a cookie based on session
			cookie := http.Cookie{
				Name:     cookieName,
				Value:    sess.Uuid,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/", 302)
		} else {
			http.Redirect(w, r, "/login", 302)
		}
	}
}

// Check if the session from Cookie is valid
func isValidSession(r *http.Request, app *data.App) (*data.Session, bool) {
	// get the session from the cookie
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Info("Cannot find cookie from the request object")
		return nil, false
	}

	// find the session object by the uuid from the cookie
	s := &data.Session{Uuid: cookie.Value}
	err = app.FindSessionByUuid(s)
	if err != nil {
		// there was a problem in finding the session, hence invalid session
		// Log the error and return accordingly
		log.Warn("Cannot find session by uuid:", cookie.Value)
		return nil, false
	} else {
		return s, true
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	t, err := parseTemplateFiles("login.layout.html", "public.navbar.html", "login.html")
	if err != nil {
		writeErrorToClient("Some error occurred. Please try later.", w, http.StatusInternalServerError)
		log.Error(err)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		writeErrorToClient("Some error occurred. Please try later.", w, http.StatusInternalServerError)
		log.Error(err)
	}
}

func logout(app *data.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// invalidate session
		// then redirect to "/"
		if sess, ok := isValidSession(r, app); ok {
			// so invalidate the session
			err := app.DeleteSession(sess)
			if err != nil {
				log.Error("Error deleting session from db", err)
			}
		}
		http.Redirect(w, r, "/", 302)
	}
}
