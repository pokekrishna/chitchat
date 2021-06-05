package httpserver

import (
	"context"
	"fmt"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
	"reflect"
	"runtime"
)

func writeErrorToClient(message string, w http.ResponseWriter){
	if _, err := w.Write([]byte(message)); err != nil {
		log.Error("cannot write error message to client:", message)
	}
}


// logHandler chains the called upon handler adds the logging
// logic before returning an anonymous function.
func logHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// do logging and then call handler
		name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
		log.Info(fmt.Sprintf("Handler function called - %s", name))
		handler(w, r)
	}
}

// Shutdown initiates graceful shutdown and can help perform clean up tasks before server shutdown
func Shutdown(ctx context.Context, server * http.Server) error {
	log.Info("Shutting down Server...")

	// FUTURE: perform clean up tasks here

	if err := server.Shutdown(ctx); err != nil{
		// Error from closing listeners or context error
		log.Error("HTTP server shutdown", err)
		return err
	}
	return nil
}
