package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/city"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
)

// GetOffering gets the offering of WoningFinder (plans, cities and housing type)
func (h *handler) GetOffering(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Plan                 []customer.Plan `json:"plan"`
		SupportedCities      []*city.City    `json:"supported_cities"`
		SupportedHousingType []string        `json:"supported_housing_type"`
	}

	var offering response

	// get cities
	cities, err := h.corporationService.GetCities()
	if err != nil {
		errorMsg := fmt.Errorf("error while getting offering")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}
	offering.SupportedCities = cities

	// add supported types
	offering.SupportedHousingType = append(offering.SupportedHousingType, []string{string(corporation.HousingTypeAppartement), string(corporation.HousingTypeHouse)}...)

	// add supported plans
	offering.Plan = append(offering.Plan, []customer.Plan{customer.PlanBasis, customer.PlanPro}...)

	// return response
	json.NewEncoder(w).Encode(offering)
}
