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

// TODO: Write tests for Router. Changing logHandler needed a lot of integration testing.

func Router(ctx context.Context, db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	app := &data.App{DB: db}

	// WebApp routes
	appRouter := router.PathPrefix("/").Subrouter()

	// The following directory is relative to the location where the
	// main program is being run from.
	files := http.FileServer(http.Dir("internal/httpserver/static"))
	staticHandler := http.StripPrefix("/static/", files)
	indexHandler := index(app)
	errHandlerHandler := errHandler(app)
	loginHandler := http.HandlerFunc(login)
	logoutHandler := logout(app)
	authenticateHandler := authenticate(app)

	appRouter.Use(middleware.LoggingMiddleware)
	appRouter.PathPrefix("/static/").Handler(staticHandler.(http.HandlerFunc)).Methods(http.MethodGet)
	appRouter.Handle("/", indexHandler).Methods(http.MethodGet)
	appRouter.Handle("/err", errHandlerHandler).Methods(http.MethodGet)
	appRouter.Handle("/login", loginHandler).Methods(http.MethodGet)
	appRouter.Handle("/logout", logoutHandler).Methods(http.MethodGet)
	appRouter.Handle("/authenticate", authenticateHandler).Methods(http.MethodPost)

	// TODO : Try using go-swagger for routes and server stubs

	// API routes
	apiV1Router := router.PathPrefix("/api/v1").Subrouter()
	apiV1Router.Use(middleware.CheckRequestHeadersMiddleware(ctx))
	apiV1Router.Use(middleware.AddResponseHeadersMiddleware)
	apiV1Router.Use(middleware.LoggingMiddleware)
	apiV1Router.Handle("/threads", middleware.NewHandler(app, apiV1.Threads)).Methods(http.MethodGet)

	return router
}