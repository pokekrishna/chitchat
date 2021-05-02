package httpserver

import (
	"fmt"
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
		// user authenticated
		_, err = w.Write([]byte (fmt.Sprintf("Welcome, %s!", u.Name)))
		if err != nil {
			log.Error("cannot write response", err)
		}
	}


}