package middleware

import (
	"context"
	"net/http"
)

type Handler struct{
	ctx context.Context
	handler http.Handler
}

func (h *Handler) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.handler.ServeHTTP()
}

func NewHandler(ctx context.Context, handler http.Handler) *Handler{
	return &Handler{
		ctx: ctx,
		handler: handler,
	}
}