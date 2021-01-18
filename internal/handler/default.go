package handler

import (
	"net/http"

	"github.com/go-chi/render"
)

func (h *handler) NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	render.Render(w, r, ErrNotFound)
}

func (h *handler) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
	render.Render(w, r, ErrMethodNotAllowed)
}
