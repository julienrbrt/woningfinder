package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
)

// UserInfo gets all the user information
func (h *handler) UserInfo(w http.ResponseWriter, r *http.Request) {
	// extract jwt
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		render.Render(w, r, handlerErrors.ErrBadRequest)
		return
	}

	// get user from jwt claims
	userFromJTW, err := auth.ExtractUserFromJWT(claims)
	if err != nil {
		render.Render(w, r, handlerErrors.ErrBadRequest)
		return
	}

	user, err := h.userService.GetUser(userFromJTW)
	if err != nil {
		errorMsg := fmt.Errorf("failed get user information")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	// activate account if first time login
	if !user.Plan.IsValid() {
		plan, _ := customer.PlanFromName(user.Plan.PlanName)
		if _, err := h.userService.ConfirmUser(user.Email, plan); err != nil {
			errorMsg := fmt.Errorf("error while starting free trial (validating user)")
			h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
			render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
			return
		}
	}

	json.NewEncoder(w).Encode(user)
}
