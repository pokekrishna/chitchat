package httpserver

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/pokekrishna/chitchat/internal/data"
	"net/http"
	"github.com/pokekrishna/chitchat/internal/httpserver/api/v1"
)

func Router(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// The following directory is relative to the location where the
	// main program is being run from.
	files := http.FileServer(http.Dir("internal/httpserver/static"))
	staticHandler := http.StripPrefix("/static/", files)
	router.PathPrefix("/static/").HandlerFunc(logHandler(staticHandler.(http.HandlerFunc))).Methods(http.MethodGet)

	app := &data.App{DB: db}
	
	indexHandler := logHandler(index(app))
	errHandlerHandler := logHandler(errHandler(app))
	loginHandler := logHandler(login)
	logoutHandler := logHandler(logout(app))
	authenticateHandler := logHandler(authenticate(app))

	router.HandleFunc("/", indexHandler).Methods(http.MethodGet)
	router.HandleFunc("/err", errHandlerHandler).Methods(http.MethodGet)
	router.HandleFunc("/login", loginHandler).Methods(http.MethodGet)
	router.HandleFunc("/logout", logoutHandler).Methods(http.MethodGet)
	router.HandleFunc("/authenticate", authenticateHandler).Methods(http.MethodPost)

	// TODO : Try using go-swagger for routes and server stubs
	// API routes
	threadsHandler := logHandler(v1.Threads(app))
	router.HandleFunc("/api/v1/threads", threadsHandler).Methods(http.MethodGet)

	return router
}
