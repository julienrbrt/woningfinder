package handler

import (
	_ "embed"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	handlerErrors "github.com/julienrbrt/woningfinder/internal/handler/errors"
	"github.com/julienrbrt/woningfinder/pkg/util"
	"go.uber.org/zap"
)

type contactFormRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Message  string `json:"message"`
	AntiSpam int    `json:"phone"` // using phone to lure bots
}

// Bind permits go-chi router to verify the user input and marshal it
func (c *contactFormRequest) Bind(r *http.Request) error {
	if c.Email == "" || c.Name == "" || c.Message == "" || c.AntiSpam == 0 {
		return errors.New("fields cannot be empty")
	}

	if !util.IsEmailValid(c.Email) {
		return fmt.Errorf("email invalid")
	}

	// verify anti spam
	// the anti spam is build the sum of all bytes of the sent email and message + 374
	var sum = 374
	bytes := append([]byte(c.Email), []byte(c.Message)...)
	for _, b := range bytes {
		sum += int(b)
	}

	if c.AntiSpam != sum {
		return errors.New("contact form cannot be processed")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*contactFormRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) ContactForm(w http.ResponseWriter, r *http.Request) {
	message := &contactFormRequest{}
	if err := render.Bind(r, message); err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(err))
		return
	}

	// send contact message to woningfinder
	if err := h.emailService.ContactFormSubmission(message.Name, message.Email, message.Message); err != nil {
		errorMsg := "failed sending message: please try again"
		h.logger.Error(errorMsg, zap.Error(err))
		render.Render(w, r, handlerErrors.ServerErrorRenderer(fmt.Errorf(errorMsg)))
		return
	}

	// returns 200 by default
}
