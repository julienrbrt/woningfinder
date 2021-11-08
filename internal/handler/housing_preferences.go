package handler

import (
	"fmt"
	"net/http"

	jwtauth "github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	"go.uber.org/zap"
)

type updateHousingPreferencesRequest struct {
	*customer.HousingPreferences
}

// Bind permits go-chi router to verify the user input and marshal it
func (u *updateHousingPreferencesRequest) Bind(r *http.Request) error {
	return u.HasMinimal()
}

// Render permits go-chi router to render the user
func (*updateHousingPreferencesRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// UpdateUserInfo updates user housing preferences
func (h *handler) UpdateHousingPreferences(w http.ResponseWriter, r *http.Request) {
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

	request := &updateHousingPreferencesRequest{}
	if err := render.Bind(r, request); err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(err))
		return
	}

	// update housing preferences
	if err := h.userService.UpdateHousingPreferences(userFromJWT.ID, request.HousingPreferences); err != nil {
		errorMsg := "failed to update housing information"
		h.logger.Error(errorMsg, zap.Error(err))
		render.Render(w, r, handlerErrors.ServerErrorRenderer(fmt.Errorf(errorMsg)))
		return
	}

	// returns 200 by default
}
