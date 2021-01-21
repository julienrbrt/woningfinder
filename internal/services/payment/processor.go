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

	if err := s.redisClient.Publish(pubSubPayment, result); err != nil {
		return fmt.Errorf("error loading payment %v into queue: %w", payment, err)
	}

	return nil
}

// ProcessPayment read the queued payment from
func (s *service) ProcessPayment(paymentCh chan<- entity.PaymentData) error {
	ch, err := s.redisClient.Subscribe(pubSubPayment)
	if err != nil {
		return err
	}

	// Consume messages
	for msg := range ch {
		var payment entity.PaymentData
		err := json.Unmarshal([]byte(msg.Payload), &payment)
		if err != nil {
			s.logger.Sugar().Errorf("error while unmarshaling payment data: %w", err)
			continue
		}

		// send payment to channel
		paymentCh <- payment
	}

	// should never happen
	return nil
}
