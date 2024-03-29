package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	handlerErrors "github.com/julienrbrt/woningfinder/internal/handler/errors"
	"github.com/julienrbrt/woningfinder/pkg/util"
	"go.uber.org/zap"
)

type loginRequest struct {
	Email string `json:"email"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (l *loginRequest) Bind(r *http.Request) error {
	if !util.IsEmailValid(l.Email) {
		return errors.New("please give a valid email")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*loginRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Login contains the handler for sending a user one time token
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	login := &loginRequest{}
	if err := render.Bind(r, login); err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(err))
		return
	}

	user, err := h.userService.GetUser(login.Email)
	if err != nil {
		render.Render(w, r, handlerErrors.ErrNotFound)
		return
	}

	// send login email
	if err := h.emailService.SendLogin(user); err != nil {
		errorMsg := "error while sending login email"
		h.logger.Error(errorMsg, zap.Error(err))
		render.Render(w, r, handlerErrors.ServerErrorRenderer(fmt.Errorf(errorMsg)))
	}

	// returns 200 by default
}
