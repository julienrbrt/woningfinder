package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/woningfinder/woningfinder/internal/customer"
	handlerErrors "github.com/woningfinder/woningfinder/internal/handler/errors"
	"github.com/woningfinder/woningfinder/pkg/util"
)

const (
	successURL = "https://woningfinder.nl/login?thanks=true"
	cancelURL  = "https://woningfinder.nl/start/voltooien?cancelled=true"
)

type subscriptionRequest struct {
	Email string `json:"email"`
}

type subscriptionResponse struct {
	StripeSessionID string `json:"stripe_session_id"`
}

// Bind permits go-chi router to verify the user input and marshal it
func (p *subscriptionRequest) Bind(r *http.Request) error {
	if !util.IsEmailValid(p.Email) {
		return errors.New("please enter a valid email")
	}

	return nil
}

// Render permits go-chi router to render the user
func (*subscriptionRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *handler) Subscription(w http.ResponseWriter, r *http.Request) {
	request := &subscriptionRequest{}
	if err := render.Bind(r, request); err != nil {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(err))
		return
	}

	// get given user
	user, err := h.userService.GetUser(request.Email)
	if err != nil {
		render.Render(w, r, handlerErrors.ErrNotFound)
		return
	}

	if user.Plan.IsFree() {
		render.Render(w, r, handlerErrors.BadRequestErrorRenderer(errors.New("user has free plan")))
		return
	}

	plan, err := customer.PlanFromName(user.Plan.Name)
	if err != nil {
		render.Render(w, r, handlerErrors.ServerErrorRenderer(err))
		return
	}

	// process payment by creating a Stripe session ID
	h.createStripeCheckoutSession(request.Email, plan, w, r)
}
