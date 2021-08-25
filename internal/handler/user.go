package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	userFromJWT, err := auth.ExtractUserFromJWT(claims)
	if err != nil {
		render.Render(w, r, handlerErrors.ErrBadRequest)
		return
	}

	user, err := h.userService.GetUser(userFromJWT)
	if err != nil {
		errorMsg := fmt.Errorf("failed get user information")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	if !user.Plan.IsValid() {
		if user.Plan.FreeTrialStartedAt == (time.Time{}) { // user is invalid and no start free trial means not activated user
			if _, err := h.userService.ConfirmUser(user.Email); err != nil {
				errorMsg := fmt.Errorf("error while starting free trial (validating user)")
				h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
				render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
				return
			}

		} else if !user.Plan.IsFreeTrialValid() { // user is invalid with free trial started mean free trial expired user
			if err := h.emailService.SendFreeTrialReminder(user); err != nil {
				errorMsg := fmt.Errorf("error while sending free trial reminder")
				h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
				render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
				return
			}
		}
	}

	json.NewEncoder(w).Encode(struct {
		*customer.User
		ValidPlan bool `json:"valid_plan"`
	}{
		user, user.Plan.IsValid(),
	})
}
