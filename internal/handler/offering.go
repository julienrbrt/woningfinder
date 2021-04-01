package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	handlerEntity "github.com/woningfinder/woningfinder/internal/handler/entity"
)

// GetOffering gets the offering of WoningFinder (plans, cities and housing type)
func (h *handler) GetOffering(w http.ResponseWriter, r *http.Request) {
	type plan struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}

	type response struct {
		Plan                 []plan   `json:"plan"`
		SupportedCities      []string `json:"supported_cities"`
		SupportedHousingType []string `json:"supported_housing_type"`
	}

	var offering response

	// get cities
	cities, err := h.corporationService.GetCities()
	if err != nil {
		render.Render(w, r, handlerEntity.ServerErrorRenderer(fmt.Errorf("error while getting offering")))
		return
	}

	for _, c := range *cities {
		offering.SupportedCities = append(offering.SupportedCities, c.Name)
	}

	// add supported types
	offering.SupportedHousingType = append(offering.SupportedHousingType, []string{string(entity.HousingTypeAppartement), string(entity.HousingTypeHouse)}...)

	// add supported plans
	offering.Plan = append(offering.Plan, plan{Name: string(entity.PlanBasis), Price: entity.PlanBasis.Price()})
	offering.Plan = append(offering.Plan, plan{Name: string(entity.PlanPro), Price: entity.PlanPro.Price()})

	// return response
	json.NewEncoder(w).Encode(offering)
}