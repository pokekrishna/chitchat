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
	router.PathPrefix("/static/").HandlerFunc(logHandler(staticHandler.(http.HandlerFunc)))

	router.HandleFunc("/", logHandler(index))
	router.HandleFunc("/err", logHandler(errHandler))
	router.HandleFunc("/login", logHandler(login))
	router.HandleFunc("/logout", logHandler(logout))
	router.HandleFunc("/authenticate", logHandler(authenticate))

	return router
}
