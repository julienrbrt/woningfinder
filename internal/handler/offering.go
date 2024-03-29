package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	handlerErrors "github.com/julienrbrt/woningfinder/internal/handler/errors"
	"go.uber.org/zap"
)

// GetOffering gets the offering of WoningFinder (cities and housing type)
func (h *handler) GetOffering(w http.ResponseWriter, r *http.Request) {
	type response struct {
		SupportedCities       []*city.City `json:"supported_cities"`
		SupportedHousingTypes []string     `json:"supported_housing_types"`
	}

	var offering response

	// get cities
	cities, err := h.corporationService.GetCities()
	if err != nil {
		errorMsg := "error while getting offering"
		h.logger.Error(errorMsg, zap.Error(err))
		render.Render(w, r, handlerErrors.ServerErrorRenderer(fmt.Errorf(errorMsg)))
		return
	}
	offering.SupportedCities = cities

	// add supported types
	offering.SupportedHousingTypes = append(offering.SupportedHousingTypes, []string{string(corporation.HousingTypeAppartement), string(corporation.HousingTypeHouse)}...)

	// return response
	json.NewEncoder(w).Encode(offering)
}
