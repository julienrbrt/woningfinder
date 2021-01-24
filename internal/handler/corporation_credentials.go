package handler

import (
	"encoding/json"
	"net/http"

	"github.com/woningfinder/woningfinder/internal/auth"

	"github.com/go-chi/jwtauth"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	handlerEntity "github.com/woningfinder/woningfinder/internal/handler/entity"
)

// GetCorporationCredentials gets a list of corporation credentials that match the user housing preferences
func (h *handler) GetCorporationCredentials(w http.ResponseWriter, r *http.Request) {
	// extract jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		render.Render(w, r, handlerEntity.ErrBadRequest)
		return
	}

	// get user from jwt claims
	user, err := auth.ExtractUserFromJWT(claims)
	if err != nil {
		render.Render(w, r, handlerEntity.ErrBadRequest)
		return
	}

	corporations, err := h.userService.GetHousingPreferencesMatchingCorporation(user)
	if err != nil {
		render.Render(w, r, handlerEntity.ServerErrorRenderer(err))
		return
	}

	// TODO gets if the credentials are stored or not

	var credentials []handlerEntity.CredentialsResponse
	for _, corporation := range corporations {
		credentials = append(credentials, handlerEntity.CredentialsResponse{CorporationName: corporation.Name})
	}

	json.NewEncoder(w).Encode(credentials)
}

// UpdateCorporationCredentials permits to update the given corporation credentials of an user
func (h *handler) UpdateCorporationCredentials(w http.ResponseWriter, r *http.Request) {
	// extract jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		render.Render(w, r, handlerEntity.ErrBadRequest)
		return
	}

	// get user from jwt claims
	user, err := auth.ExtractUserFromJWT(claims)
	if err != nil {
		render.Render(w, r, handlerEntity.ErrBadRequest)
		return
	}

	var credentials handlerEntity.CredentialsRequest
	if err := render.Bind(r, &credentials); err != nil {
		render.Render(w, r, handlerEntity.ErrBadRequest)
		return
	}

	corporationCredentials := entity.CorporationCredentials{
		UserID:          user.ID,
		CorporationName: credentials.CorporationName,
		Corporation:     entity.Corporation{Name: credentials.CorporationName},
		Login:           credentials.Login,
		Password:        credentials.Password,
	}
	if err := h.userService.CreateCorporationCredentials(user, corporationCredentials); err != nil {
		render.Render(w, r, handlerEntity.ServerErrorRenderer(err))
		return
	}

	// returns 200 by default
}

// DeleteCorporationCredentials permits to delete a corporation credentials
func (h *handler) DeleteCorporationCredentials(w http.ResponseWriter, r *http.Request) {
	// TODO implement DeleteCorporationCredentials
}
