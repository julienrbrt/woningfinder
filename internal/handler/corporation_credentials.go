package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/entity"
)

// GetCorporationCredentials gets a list of corporation credentials that match the user housing preferences
func (h *handler) GetCorporationCredentials(w http.ResponseWriter, r *http.Request) {
	// extract jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		render.Render(w, r, entity.ErrBadRequest)
		return
	}

	// get user from jwt claims
	user, err := auth.ExtractUserFromJWT(claims)
	if err != nil {
		render.Render(w, r, entity.ErrBadRequest)
		return
	}

	corporations, err := h.userService.GetHousingPreferencesMatchingCorporation(user)
	if err != nil {
		errorMsg := fmt.Errorf("failed getting housing corporation relevant for you")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, entity.ServerErrorRenderer(errorMsg))
		return
	}

	// used to display which housing corporation are supported for the user housing preferences
	type response struct {
		CorporationName string `json:"corporation_name"`
		CorporationURL  string `json:"corporation_url"`
		IsKnown         bool   `json:"is_known"`
	}

	var credentials []response
	for _, corporation := range corporations {
		var isKnown bool
		if creds, err := h.userService.GetCorporationCredentials(user.ID, corporation.Name); err == nil {
			if creds.Login != "" {
				isKnown = true
			}
		}

		credentials = append(credentials, response{
			CorporationName: corporation.Name,
			CorporationURL:  corporation.URL,
			IsKnown:         isKnown,
		})
	}

	json.NewEncoder(w).Encode(credentials)
}

// UpdateCorporationCredentials permits to update the given corporation credentials of an user
func (h *handler) UpdateCorporationCredentials(w http.ResponseWriter, r *http.Request) {
	// extract jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		render.Render(w, r, entity.ErrBadRequest)
		return
	}

	// get user from jwt claims
	user, err := auth.ExtractUserFromJWT(claims)
	if err != nil {
		render.Render(w, r, entity.ErrBadRequest)
		return
	}

	var credentials entity.CorporationCredentials
	if err := render.Bind(r, &credentials); err != nil {
		render.Render(w, r, entity.ErrBadRequest)
		return
	}

	corporationCredentials := entity.CorporationCredentials{
		UserID:          user.ID,
		CorporationName: credentials.CorporationName,
		Login:           credentials.Login,
		Password:        credentials.Password,
	}
	if err := h.userService.CreateCorporationCredentials(user.ID, corporationCredentials); err != nil {
		errorMsg := fmt.Errorf("failed creating corporation credentials")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, entity.ServerErrorRenderer(errorMsg))
		return
	}

	// returns 200 by default
}

// DeleteCorporationCredentials permits to delete a corporation credentials
func (h *handler) DeleteCorporationCredentials(w http.ResponseWriter, r *http.Request) {
	// TODO implement DeleteCorporationCredentials
	panic("not implemented")
}
