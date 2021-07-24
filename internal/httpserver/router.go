package httpserver

import (
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/middleware"
	"net/http"
	apiV1 "github.com/pokekrishna/chitchat/internal/httpserver/api/v1"
)

func Router(ctx context.Context, db *sql.DB) *mux.Router {
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
	// TODO: not checking req headers
	// API routes

	// TODO: convert logHandler to mux.Middlewarefunc and add it as middleware
	//threadsHandler := logHandler(apiV1.Threads(ctx, app))

	apiV1Router := router.PathPrefix("/api/v1").Subrouter()
	apiV1Router.Use(middleware.CheckRequestHeadersMiddleware(ctx))
	apiV1Router.Use(middleware.AddResponseHeadersMiddleware)
	apiV1Router.Handle("/threads", middleware.NewHandler(app, apiV1.Threads)).Methods(http.MethodGet)

	return router
}