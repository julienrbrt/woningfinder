package payment

import (
	"fmt"
	"time"

	"github.com/woningfinder/woningfinder/internal/customer"
)

// ProcessPayment validate user account
func (s *service) ProcessFreeTrial(email string, plan customer.Plan) error {
	// get user
	user, err := s.userService.GetUser(&customer.User{Email: email})
	if err != nil {
		return fmt.Errorf("error while processing payment data: cannot get user: %w", err)
	}

	// set that has started its free trial
	if _, err := s.dbClient.Conn().
		Model(&customer.UserPlan{UserID: user.ID, PlanName: plan.Name}).
		OnConflict("(user_id) DO UPDATE").
		Insert(); err != nil {
		return fmt.Errorf("error when adding user plan: %w", err)
	}

	if err := s.emailService.SendWelcome(user); err != nil {
		return fmt.Errorf("error while processing payment data: %w", err)
	}

	return nil
}

// ProcessPayment set that the user has paid
func (s *service) ProcessPayment(email string, plan customer.Plan) error {
	// get user
	user, err := s.userService.GetUser(&customer.User{Email: email})
	if err != nil {
		return fmt.Errorf("error while processing payment data: cannot get user: %w", err)
	}

	// set that user has paid
	if _, err := s.dbClient.Conn().
		Model(&customer.UserPlan{UserID: user.ID, PlanName: plan.Name, PurchasedAt: time.Now()}).
		OnConflict("(user_id) DO UPDATE").
		Insert(); err != nil {
		return fmt.Errorf("error when adding user plan: %w", err)
	}

	if err := s.emailService.SendThankYou(user); err != nil {
		return fmt.Errorf("error while processing payment data: %w", err)
	}

	return nil
}
