package matcher

import (
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/customer"
)

type Matcher interface {
	MatchOffer(user customer.User, offer corporation.Offer) bool
}

type matcher struct{}

func NewMatcher() Matcher {
	return &matcher{}
}

func (m *matcher) MatchOffer(user customer.User, offer corporation.Offer) bool {
	return m.matchCriteria(user, offer) && m.matchPreferences(user.HousingPreferences, offer)
}
