package payment

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// ProcessPayment set that the user has paid
func (s *service) ProcessPayment(payment *entity.PaymentData) error {
	// get user
	user, err := s.userService.GetUser(&entity.User{Email: payment.UserEmail})
	if err != nil {
		return fmt.Errorf("error while processing payment data: cannot get user: %w", err)
	}

	// set that user has paid
	if err := s.userService.SetPaid(user, payment.Plan); err != nil {
		return fmt.Errorf("error while processing payment data: %w", err)
	}

	if err := s.notifyUser(user); err != nil {
		return fmt.Errorf("error while processing payment data: %w", err)
	}

	return nil
}

// notifyUser sends confirmation email
func (s *service) notifyUser(user *entity.User) error {
	return s.notificationsService.SendWelcome(user)
}
