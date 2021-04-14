package matcher

import (
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/customer"
)

// MatchCriteria verifies that an user match the offer criterias
func (m *matcher) matchCriteria(user customer.User, offer corporation.Offer) bool {
	age := time.Now().Year() - user.BirthYear

	// checks if offer age is set and check boundaries
	if offer.MinAge > 0 && ((age < offer.MinAge) || (offer.MaxAge != 0 && age > offer.MaxAge)) {
		return false
	}

	// checks if offer family size is set and check boundaries
	if offer.MinFamilySize > 0 && (user.FamilySize < offer.MinFamilySize) || (offer.MaxFamilySize > 0 && user.FamilySize > offer.MaxFamilySize) {
		return false
	}

	// checks if offer incomes is set and check boundaries
	min, max := m.passendToewijzen(user)
	if offer.Housing.Price < min || offer.Housing.Price > max {
		return false
	}

	return true
}
