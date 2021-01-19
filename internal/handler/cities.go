package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

// GetCities gets the cities supported by WoningFinder
func (h *handler) GetCities(w http.ResponseWriter, r *http.Request) {
	cities, err := h.corporationService.GetCities()
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(errors.New("error while getting cities")))
		return
	}

	json.NewEncoder(w).Encode(cities)
}
