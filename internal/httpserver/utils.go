package httpserver

import (
	"context"
	"fmt"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
	"os"
	"reflect"
	"runtime"
)

// TODO: validate http status
func writeErrorToClient(message string, w http.ResponseWriter, httpStatus int) {
	if err := checkWriteHeaderCode(httpStatus); err != nil{
		log.Error(err)
		httpStatus = http.StatusInternalServerError // fallback status
	}
	w.WriteHeader(httpStatus)
	if _, err := w.Write([]byte(message)); err != nil {
		log.Error("cannot write error message to client:", message)
	}
}

// checkWriteHeaderCode is inspired by https://golang.org/src/net/http/server.go
func checkWriteHeaderCode(code int) error{
	// Issue 22880: require valid WriteHeader status codes.
	// For now we only enforce that it's three digits.
	// In the future we might block things over 599 (600 and above aren't defined
	// at https://httpwg.org/specs/rfc7231.html#status.codes)
	// and we might block under 200 (once we have more mature 1xx support).
	// But for now any three digits.
	//
	// We used to send "HTTP/1.1 000 0" on the wire in responses but there's
	// no equivalent bogus thing we can realistically send in HTTP/2,
	// so we'll consistently panic instead and help people find their bugs
	// early. (We can't return an error from WriteHeader even if we wanted to.)
	if code < 100 || code > 999 {
		return fmt.Errorf("invalid WriteHeader code %v", code)
	}
	return nil
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
func Shutdown(ctx context.Context, server *http.Server) error {
	log.Info("Shutting down Server...")

	// FUTURE: perform clean up tasks here

	if err := server.Shutdown(ctx); err != nil {
		// Error from closing listeners or context error
		log.Error("HTTP server shutdown", err)
		return err
	}
	return nil
}

// ChangeDirForTest is a convenience function for tests that changes the current
// directory by going 2 steps up. or "../../" This is required because while
// running the tests, go test command goes inside individual package dir. and
// runs the tests. This is different from when the app is run. This change in
// directory causes discrepancy in finding the template files.
//
// ChangeDirForTest returns the directory which was the directory before
// executing this function. This is required because it is important for Test
// functions to reset back the directory before completing the execution.
// Preferred way to use this is mentioned below.
//
// Usage:
//	defer os.Chdir(ChangeDirForTest())
func ChangeDirForTest() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic("Cannot get current working directory")
	}

	if err := os.Chdir("../.."); err != nil {
		panic("Cannot change directory")
	}
	return cwd
}
