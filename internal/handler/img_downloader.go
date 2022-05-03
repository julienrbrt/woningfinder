package handler

import (
	"net/http"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	handlerErrors "github.com/julienrbrt/woningfinder/internal/handler/errors"
)

// GetOfferImage gets the image of an offer
func (h *handler) GetOfferImage(w http.ResponseWriter, r *http.Request) {
	img, err := h.imgClient.Get(chi.URLParam(r, "img"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		render.Render(w, r, handlerErrors.ErrNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(img)
}
