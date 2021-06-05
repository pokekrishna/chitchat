package httpserver

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/pokekrishna/chitchat/internal/data"
	"net/http"
)

func Router(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// The following directory is relative to the location where the
	// main program is being run from.
	files := http.FileServer(http.Dir("internal/httpserver/static"))
	staticHandler:= http.StripPrefix("/static/", files)
	router.PathPrefix("/static/").HandlerFunc(logHandler(staticHandler.(http.HandlerFunc))).Methods(http.MethodGet)

	// TODO: is it safe to initialize models here and reuse ...
	// TODO: ... from concurrency standpoint.
	t := data.NewThread(db)
	u := data.NewUser(db)
	s := data.NewSession(db, u)

	indexHandler := logHandler(index(t))
	errHandlerHandler := logHandler(errHandler(s))
	loginHandler := logHandler(login)
	logoutHandler := logHandler(logout(s))
	authenticateHandler := logHandler(authenticate(u))

	router.HandleFunc("/", indexHandler).Methods(http.MethodGet)
	router.HandleFunc("/err", errHandlerHandler).Methods(http.MethodGet)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	router.HandleFunc("/logout", logoutHandler).Methods(http.MethodGet)
	router.HandleFunc("/authenticate", authenticateHandler).Methods(http.MethodPost)

	return router
}
