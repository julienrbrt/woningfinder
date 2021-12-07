package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwtauth "github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"github.com/julienrbrt/woningfinder/internal/auth"
	"github.com/julienrbrt/woningfinder/internal/customer"
	handlerErrors "github.com/julienrbrt/woningfinder/internal/handler/errors"
	"go.uber.org/zap"
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
		errorMsg := "failed to get user information"
		h.logger.Error(errorMsg, zap.Error(err))
		render.Render(w, r, handlerErrors.ServerErrorRenderer(fmt.Errorf(errorMsg)))
		return
	}

	// confirm user for first login
	if !user.IsActivated() {
		if err := h.userService.ConfirmUser(user.Email); err != nil {
			errorMsg := "error while activating user"
			h.logger.Error(errorMsg, zap.Error(err))
			render.Render(w, r, handlerErrors.ServerErrorRenderer(fmt.Errorf(errorMsg)))
			return
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
		TotalReaction int `json:"total_reaction"`
	}{
		user,
		totalReaction,
	})
}

type updateUserInfoRequest struct {
	Name             string `json:"name"`
	YearlyIncome     int    `json:"yearly_income"`
	FamilySize       int    `json:"family_size"`
	HasAlertsEnabled bool   `json:"has_alerts_enabled"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (u *updateUserInfoRequest) Bind(r *http.Request) error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}

	if u.YearlyIncome < 0 {
		return fmt.Errorf("user yearly income invalid")
	}

	if u.FamilySize < 0 {
		return fmt.Errorf("user family size invalid")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*updateUserInfoRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// UpdateUserInfo updates the user basic information
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

	request := &updateUserInfoRequest{}
	if err := render.Bind(r, request); err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(err))
		return
	}

	// update user
	userFromJWT.Name = request.Name
	userFromJWT.YearlyIncome = request.YearlyIncome
	userFromJWT.FamilySize = request.FamilySize
	userFromJWT.HasAlertsEnabled = request.HasAlertsEnabled

	if err := h.userService.UpdateUser(userFromJWT); err != nil {
		errorMsg := "failed to update user information"
		h.logger.Error(errorMsg, zap.Error(err))
		render.Render(w, r, handlerErrors.ServerErrorRenderer(fmt.Errorf(errorMsg)))
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
			h.logger.Error("failed to send email", zap.Error(err))
		}
	}

	if len(request.Feedback) > 0 {
		if err := h.emailService.ContactFormSubmission("Deleted user", user.Email, request.Feedback); err != nil {
			h.logger.Error("failed to send email", zap.Error(err))
		}
	}

	errorMsg := "failed to delete user"

	// get user
	user, err = h.userService.GetUser(user.Email)
	if err != nil {
		h.logger.Error(errorMsg, zap.Error(err))
		render.Render(w, r, handlerErrors.ServerErrorRenderer(fmt.Errorf(errorMsg)))
		return
	}

	// delete user
	if err := h.userService.DeleteUser(user.Email); err != nil {
		h.logger.Error(errorMsg, zap.Error(err))
		render.Render(w, r, handlerErrors.ServerErrorRenderer(fmt.Errorf(errorMsg)))
		return
	}

}
