package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	handlerEntity "github.com/woningfinder/woningfinder/internal/handler/entity"
)

// GetCities gets the cities supported by WoningFinder
func (h *handler) GetCities(w http.ResponseWriter, r *http.Request) {
	cities, err := h.corporationService.GetCities()
	if err != nil {
		render.Render(w, r, handlerEntity.ServerErrorRenderer(fmt.Errorf("error while getting cities")))
		return
	}

	json.NewEncoder(w).Encode(cities)
}
