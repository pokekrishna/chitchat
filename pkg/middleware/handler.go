package middleware

import (
	"github.com/pokekrishna/chitchat/internal/data"
	"net/http"
)
// TODO: complete docs of the types and methods of this file

// WorkerFunc is the func that does request handling but does not directly have
// method ServeHTTP
type WorkerFunc func(*data.App, http.ResponseWriter, *http.Request)

type Handler struct{
	app *data.App
	hf  WorkerFunc
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.hf(h.app, w, r)
}

// Need this constructor so that a WorkerFunc can be adapted to http.Handler
func NewHandler(app *data.App, hf WorkerFunc) *Handler{
	return &Handler{
		app: app,
		hf: hf,
	}
}