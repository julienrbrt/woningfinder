package user

import (
	"time"

	"github.com/woningfinder/woningfinder/internal/corporation"
)

// MatchCriteria verifies that an user match the offer criterias
func (u *User) MatchCriteria(offer corporation.Offer) bool {
	age := time.Now().Year() - u.BirthYear
	// checks if offer age is set and check boundaries
	if offer.MinAge > 0 && ((age < offer.MinAge) || (offer.MaxAge != 0 && age > offer.MaxAge)) {
		return false
	}

	// checks if offer family size is set and check boundaries
	if offer.MinFamilySize > 0 && (u.FamilySize < offer.MinFamilySize || (offer.MaxFamilySize > 0 && u.FamilySize > offer.MaxFamilySize)) {
		return false
	}

	// checks if offer incomes is set and check boundaries
	if offer.MinIncome > 0 && u.YearlyIncome > -1 && (u.YearlyIncome < offer.MinIncome || (offer.MaxIncome > 0 && u.YearlyIncome > offer.MaxIncome)) {
		return false
	}

	return true
}
