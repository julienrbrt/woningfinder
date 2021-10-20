package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	userService "github.com/woningfinder/woningfinder/internal/services/user"
)

type registerRequest struct {
	*customer.User
}

// Bind permits go-chi router to verify the user input and marshal it
func (u *registerRequest) Bind(r *http.Request) error {
	return u.HasMinimal()
}

// Render permits go-chi router to render the user
func (*registerRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Register contains the handler for registering on WoningFinder
func (h *handler) Register(w http.ResponseWriter, r *http.Request) {
	user := &registerRequest{}
	if err := render.Bind(r, user); err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(err))
		return
	}

	// lowercase email
	user.Email = strings.ToLower(user.Email)

	if err := h.userService.CreateUser(user.User); err != nil {
		errorMsg := fmt.Errorf("error while creating user")

		if errors.Is(err, userService.ErrUserAlreadyExist) {
			render.Render(w, r, handlerErrors.BadRequestErrorRenderer(fmt.Errorf("%s: %s", errorMsg, err.Error())))
			return
		}

		h.logger.Sugar().Errorf("%s: %w", errorMsg, err)
		render.Render(w, r, handlerErrors.ServerErrorRenderer(errorMsg))
		return
	}

	// send activation email
	if err := h.emailService.SendActivationEmail(user.User); err != nil {
		// just logging error
		h.logger.Sugar().Errorf("error while sending activation email: %w", err)
	}
}
