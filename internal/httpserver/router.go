package httpserver

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	// The following directory is relative to the location where the
	// main program is being run from.
	files := http.FileServer(http.Dir("internal/httpserver/static"))
	staticHandler:= http.StripPrefix("/static/", files)
	router.PathPrefix("/static/").HandlerFunc(logHandler(staticHandler.(http.HandlerFunc))).Methods(http.MethodGet)

	router.HandleFunc("/", logHandler(index)).Methods(http.MethodGet)
	router.HandleFunc("/err", logHandler(errHandler)).Methods(http.MethodGet)
	router.HandleFunc("/login", logHandler(login)).Methods(http.MethodGet)
	router.HandleFunc("/logout", logHandler(logout)).Methods(http.MethodGet)
	router.HandleFunc("/authenticate", logHandler(authenticate)).Methods(http.MethodPost)

	return router
}
