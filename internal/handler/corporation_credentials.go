package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// GetCorporationCredentials gets a list of corporation credentials that match the user housing preferences
func (h *handler) GetCorporationCredentials(w http.ResponseWriter, r *http.Request) {
	// TODO check JWT token and get user with userID
	var user entity.User

	corporations, err := h.userService.GetHousingPreferencesMatchingCorporation(&user)
	if err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}

	// TODO  gets if the credentials are stored or not

	var credentials []credentialsResponse
	for _, corporation := range corporations {
		credentials = append(credentials, credentialsResponse{CorporationName: corporation.Name})
	}

	json.NewEncoder(w).Encode(credentials)
}

func (h *handler) UpdateCorporationCredentials(w http.ResponseWriter, r *http.Request) {
	// TODO check JWT token and get user with userID
	var user entity.User

	var credentials credentialsRequest
	if err := render.Bind(r, &credentials); err != nil {
		render.Render(w, r, ErrorRenderer(err))
		return
	}

	corporationCredentials := entity.CorporationCredentials{
		UserID:          user.ID,
		CorporationName: credentials.CorporationName,
		Corporation:     entity.Corporation{Name: credentials.CorporationName},
		Login:           credentials.Login,
		Password:        credentials.Password,
	}
	if err := h.userService.CreateCorporationCredentials(&user, corporationCredentials); err != nil {
		render.Render(w, r, ServerErrorRenderer(err))
		return
	}

}

func (h *handler) DeleteCorporationCredentials(w http.ResponseWriter, r *http.Request) {

}
