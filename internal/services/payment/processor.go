package payment

import (
	"encoding/json"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// QueuePayment should be triggered by the payment webhook. It loads the user payment data into the payment queue
func (s *service) QueuePayment(payment *entity.PaymentData) error {
	result, err := json.Marshal(payment)
	if err != nil {
		return fmt.Errorf("error while marshaling payment for %s: %w", payment.UserEmail, err)
	}

	if err := s.redisClient.Push(paymentQueue, result); err != nil {
		return fmt.Errorf("error loading payment %v into queue: %w", payment, err)
	}

	return nil
}

// ProcessPayment read the payment queue
func (s *service) ProcessPayment() error {
	for {
		payments, err := s.redisClient.BPop(paymentQueue)
		if err != nil {
			return err
		}

		// Consume payments from queue
		for _, pay := range payments {
			var payment entity.PaymentData
			if err := json.Unmarshal([]byte(pay), &payment); err != nil {
				s.logger.Sugar().Errorf("error while unmarshaling payment data: %w", err)
				continue
			}

			// verify payment
			if err := s.hasPaid(payment); err != nil {
				s.logger.Sugar().Error(err)
			}
		}
	}
}

func (s *service) hasPaid(payment entity.PaymentData) error {
	// get user
	user, err := s.userService.GetUser(&entity.User{Email: payment.UserEmail})
	if err != nil {
		return fmt.Errorf("error while processing payment data: cannot get user: %w", err)
	}

	// set that user has paid
	if err := s.userService.SetPaid(user, payment.Plan); err != nil {
		return fmt.Errorf("error while processing payment data: %w", err)
	}

	// send confirmation email
	if err := s.notificationsService.SendWelcome(user); err != nil {
		return fmt.Errorf("error while processing payment data: %w", err)
	}

	return nil
}
