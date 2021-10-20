package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	jwtauth "github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
)

type corporationCredentialsRequest struct {
	*customer.CorporationCredentials
}

// Bind permits go-chi router to verify the user input and marshal it
func (c *corporationCredentialsRequest) Bind(r *http.Request) error {
	if c.CorporationName == "" || c.Login == "" || c.Password == "" {
		return errors.New("credentials cannot be empty")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*corporationCredentialsRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// GetCorporationCredentials gets a list of corporation credentials that match the user housing preferences
func (h *handler) GetCorporationCredentials(w http.ResponseWriter, r *http.Request) {
	// extract jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		render.Render(w, r, handlerErrors.ErrBadRequest)
		return
	}

	// get user from jwt claims
	userFromJWT, err := auth.ExtractUserFromJWT(claims)
	if err != nil {
		render.Render(w, r, handlerErrors.ErrBadRequest)
		return
	}

	corporations, err := h.userService.GetHousingPreferencesMatchingCorporation(userFromJWT.ID)
	if err != nil {
		errorMsg := fmt.Errorf("failed getting housing corporation relevant for you")
		h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
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
		if creds, err := h.userService.GetCorporationCredentials(userFromJWT.ID, corporation.Name); err == nil {
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
		render.Render(w, r, handlerErrors.ErrBadRequest)
		return
	}

	// get user from jwt claims
	user, err := auth.ExtractUserFromJWT(claims)
	if err != nil {
		render.Render(w, r, handlerErrors.ErrBadRequest)
		return
	}

	credentials := &corporationCredentialsRequest{}
	if err := render.Bind(r, credentials); err != nil {
		render.Render(w, r, handlerErrors.ErrBadRequest)
		return
	}

	corporationCredentials := &customer.CorporationCredentials{
		UserID:          user.ID,
		CorporationName: credentials.CorporationName,
		Login:           credentials.Login,
		Password:        credentials.Password,
	}

	hasCorproationCredentials, err := h.userService.HasCorporationCredentials(user.ID)
	if err != nil {
		h.logger.Sugar().Errorf("failed to get corproation credentials count: %w", err)
	}

	if err := h.userService.CreateCorporationCredentials(user.ID, corporationCredentials); err != nil {
		errorMsg := fmt.Errorf("failed creating corporation credentials")
		h.logger.Sugar().Errorf("%w: %w", errorMsg, err)

		if strings.Contains(err.Error(), userService.ErrValidationCorporationCredentials) {
			render.Render(w, r, handlerErrors.ErrUnauthorized)
		} else {
			render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		}

		return
	}

	// first time adding corporation credentials: send welcome email
	if !hasCorproationCredentials {
		if err := h.emailService.SendCorporationCredentialsFirstTimeAdded(user); err != nil {
			h.logger.Sugar().Error(err)
		}
	}

	// returns 200 by default
}
