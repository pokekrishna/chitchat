package middleware

import (
	"fmt"
	"github.com/pokekrishna/chitchat/internal/data"
	"github.com/pokekrishna/chitchat/pkg/log"
	"net/http"
	"reflect"
	"runtime"
)
// TODO: complete docs of the types and methods of this file

// WorkerFunc is the func that does request handling but does not directly have
// method ServeHTTP
type WorkerFunc func(*data.App, http.ResponseWriter, *http.Request)

type Handler struct{
	app *data.App
	wf  WorkerFunc
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.wf(h.app, w, r)
}

// Need this constructor so that a WorkerFunc can be adapted to http.Handler
func NewHandler(app *data.App, hf WorkerFunc) *Handler{
	return &Handler{
		app: app,
		wf:  hf,
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch next := next.(type) {
		case http.HandlerFunc:
			name := runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name()
			log.Info(fmt.Sprintf("Handler function called - %s", name))
		case *Handler:
			log.Info(fmt.Sprintf("Handler called - %#v", next))
		default:
			log.Error("Unknown Handler called - %v", next)
		}
		next.ServeHTTP(w, r)
	})
}