package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
	handlerEntity "github.com/woningfinder/woningfinder/internal/handler/entity"
)

// credentialsRequest defines an update and save the corporation credentials request
type credentialsRequest struct {
	CorporationName string `json:"corporation_name"`
	Login           string `json:"login"`
	Password        string `json:"password"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (c *credentialsRequest) Bind(r *http.Request) error {
	if c.CorporationName == "" || c.Login == "" || c.Password == "" {
		return errors.New("credentials cannot be empty")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*credentialsRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

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

	// used to display which housing corporation are supported for the user housing preferences
	type response struct {
		CorporationName string `json:"corporation_name"`
		IsKnown         bool   `json:"is_known"`
	}

	var credentials []response
	for _, corporation := range corporations {
		isKnown := false
		if creds, err := h.userService.GetCorporationCredentials(user.ID, corporation); err == nil {
			if creds.Login != "" {
				isKnown = true
			}
		}

		credentials = append(credentials, response{CorporationName: corporation.Name, IsKnown: isKnown})
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

	var credentials credentialsRequest
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
	if err := h.userService.CreateCorporationCredentials(user.ID, corporationCredentials); err != nil {
		render.Render(w, r, handlerEntity.ServerErrorRenderer(err))
		return
	}

	// returns 200 by default
}

// DeleteCorporationCredentials permits to delete a corporation credentials
func (h *handler) DeleteCorporationCredentials(w http.ResponseWriter, r *http.Request) {
	// TODO implement DeleteCorporationCredentials
}
