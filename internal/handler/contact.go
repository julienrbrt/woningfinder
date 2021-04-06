package handler

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
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
		return errors.New("message fields cannot be empty")
	}

	// verify anti spam
	// the anti spam is build the sum of all bytes of the sent email and message + 374
	var sum = 374
	bytes := append([]byte(c.Email), []byte(c.Message)...)
	for _, b := range bytes {
		sum += int(b)
	}

	if c.AntiSpam != sum {
		return errors.New("message cannot be sent")
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
		render.Render(w, r, entity.ErrorRenderer(err))
		return
	}

	messageTpl := `
	Hi,

	You have a new message from the WoningFinder contact form.
	
	- {{ .Name }}
	- {{ .Email }}
	
	>
	> {{ .Message }}
	>
	
	Regards,
	WoningFinder
	`

	// create body
	tpl := template.Must(template.New("contact").Parse(messageTpl))
	body := &bytes.Buffer{}
	if err := tpl.Execute(body, message); err != nil {
		errorMsg := fmt.Errorf("failed creating message: please try again")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, entity.ServerErrorRenderer(errorMsg))
		return
	}

	// send mail
	if err := h.emailClient.Send("WoningFinder Contact Submission", "", body.String(), "contact@woningfinder.nl"); err != nil {
		errorMsg := fmt.Errorf("failed sending message: please try again")
		h.logger.Sugar().Warnf("%w: %w", errorMsg, err)
		render.Render(w, r, entity.ServerErrorRenderer(errorMsg))
		return
	}

	// returns 200 by default
}
