package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwtauth "github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/auth"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
)

// GetUserInfo gets all the user information
func (h *handler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.userService.GetUser(userFromJWT.Email)
	if err != nil {
		errorMsg := fmt.Errorf("failed to get user information")
		h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	if !user.Plan.IsValid() {
		if user.Plan.FreeTrialStartedAt == (time.Time{}) { // user is invalid and no start free trial means not activated user
			if err := h.userService.ConfirmUser(user.Email); err != nil {
				errorMsg := fmt.Errorf("error while starting free trial (validating user)")
				h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
				render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
				return
			}

		} else if !user.Plan.IsFreeTrialValid() { // user is invalid with free trial started mean free trial expired user
			if err := h.emailService.SendFreeTrialReminder(user); err != nil {
				errorMsg := fmt.Errorf("error while sending free trial reminder")
				h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
				render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
				return
			}
		}
	}

	// filter housing match
	totalReaction := len(user.HousingPreferencesMatch)
	maximumHousingPreferences := 6
	if len(user.HousingPreferencesMatch) > maximumHousingPreferences {
		user.HousingPreferencesMatch = user.HousingPreferencesMatch[len(user.HousingPreferencesMatch)-maximumHousingPreferences:]
	}

	json.NewEncoder(w).Encode(struct {
		*customer.User
		ValidPlan     bool `json:"valid_plan"`
		TotalReaction int  `json:"total_reaction"`
	}{
		user,
		// consider a plan valid if user not activated for not displaying the invalid plan alert at first login of an user
		!user.Plan.IsActivated() || user.Plan.IsValid(),
		totalReaction,
	})
}

type updateUserInfoRequest struct {
	*customer.HousingPreferences
}

// Bind permits go-chi router to verify the user input and marshal it
// TODO update all well user user (right now only update housing prefernces)
func (u *updateUserInfoRequest) Bind(r *http.Request) error {
	return u.HasMinimal()
}

// Render permits go-chi router to render the user
func (*updateUserInfoRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// UpdateUserInfo updates the user and/or its housing preferences
// TODO update all well user user (right now only update housing prefernces)
func (h *handler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
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

	userInfoRequest := new(updateUserInfoRequest)
	if err := render.Bind(r, userInfoRequest); err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(err))
		return
	}

	// update housing preferences
	if err := h.userService.UpdateHousingPreferences(userFromJWT.ID, userInfoRequest.HousingPreferences); err != nil {
		errorMsg := fmt.Errorf("failed to update housing information")
		h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	// returns 200 by default
}

type deleteUserRequest struct {
	HasHouse bool   `json:"has_house"`
	Feedback string `json:"feedback"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (d *deleteUserRequest) Bind(r *http.Request) error {
	if len(d.Feedback) == 0 {
		return fmt.Errorf("feedback is required")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*deleteUserRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// DeleteUser let an user delete its account
func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
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

	var request deleteUserRequest
	if err := render.Bind(r, &request); err != nil {
		render.Render(w, r, handlerErrors.ErrBadRequest)
		return
	}

	// send feedback
	if request.HasHouse {
		if err := h.emailService.SendBye(user); err != nil {
			h.logger.Sugar().Error(err)
		}
	}

	if len(request.Feedback) > 0 {
		if err := h.emailService.ContactFormSubmission("Deleted user", user.Email, request.Feedback); err != nil {
			h.logger.Sugar().Error(err)
		}
	}

	// delete user
	if err := h.userService.DeleteUser(user.Email); err != nil {
		errorMsg := fmt.Errorf("failed to delete user")
		h.logger.Sugar().Errorf("%w: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}
}
