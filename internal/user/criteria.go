package user

import (
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// MatchCriteria verifies that an user match the offer criterias
func (u *User) MatchCriteria(offer corporation.Offer) bool {
	age := time.Now().Year() - u.BirthYear
	if (age < offer.MinAge) || (offer.MaxAge != 0 && age > offer.MaxAge) {
		return false
	}

	if u.FamilySize < offer.MinFamilySize || (offer.MaxFamilySize > 0 && u.FamilySize > offer.MaxFamilySize) {
		return false
	}

	if u.YearlyIncome > -1 && (u.YearlyIncome < offer.MinIncome || (offer.MaxIncome > 0 && u.YearlyIncome > offer.MaxIncome)) {
		return false
	}

	return true
}
