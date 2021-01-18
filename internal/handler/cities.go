package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// Cities gets the cities supported by WoningFinder
func (h *handler) Cities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	cities, err := h.GetCities()
	if err != nil {
		h.logger.Sugar().Errorf("error while getting cities: %w", err)
		render.Render(w, r, ServerErrorRenderer(errors.New("error while getting cities")))
		return
	}

	json.NewEncoder(w).Encode(cities)
}
