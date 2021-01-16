package stripe

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/domain/entity"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
)

func (c *client) CreateCustomer(u *entity.User) error {
	params := &stripe.CustomerParams{
		Description: stripe.String(fmt.Sprintf("Searching in %v\n", u.HousingPreferences[0].City)),
		Email:       stripe.String(u.Email),
		Name:        stripe.String(u.Name),
	}

	_, err := customer.New(params)
	if err != nil {
		return fmt.Errorf("error while creating stripe customer with email %s: %w", u.Email, err)
	}

	return nil
}
