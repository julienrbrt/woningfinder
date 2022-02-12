package matcher

import (
	"math"
	"time"

	"github.com/julienrbrt/woningfinder/internal/customer"
)

// last updated in april 2021

// https://www.woningmarktbeleid.nl/onderwerpen/huurtoeslag
const (
	aow   = 67
	young = 23

	MaximumIncomeSocialHouse = 43082
	Kwaliteitskortingsgrens  = 442.46
	Aftoppingsgrenslaag      = 633.25
	Aftoppingsgrenshoog      = 678.66
	Maximalehuurgrens        = 763.47
)

// PassendToewijzen determines the rent range to which an user can react
func (m *matcher) passendToewijzen(user customer.User) (float64, float64) {
	age := time.Now().Year() - user.BirthYear

	// if too rich free sector only allowed
	if user.YearlyIncome > MaximumIncomeSocialHouse {
		return Maximalehuurgrens, math.MaxInt32
	}

	// housing allowance for people younger than 23 years old and no children
	if age < young && !user.HasChildrenSameHousing {
		return 0, Kwaliteitskortingsgrens
	}

	switch {
	case user.FamilySize <= 1:
		if user.YearlyIncome <= 23725 {
			return 0, Aftoppingsgrenslaag
		} else if user.YearlyIncome <= 32200 {
			return 0, Aftoppingsgrenshoog
		} else if user.YearlyIncome <= 40024 {
			return Aftoppingsgrenshoog, Maximalehuurgrens
		}
	case user.FamilySize == 2:
		if user.YearlyIncome <= 32200 {
			return 0, Aftoppingsgrenslaag
		} else if user.YearlyIncome <= 40024 {
			return Aftoppingsgrenshoog, Maximalehuurgrens
		}
	case user.FamilySize >= 3:
		if user.YearlyIncome <= 32200 {
			return 0, Aftoppingsgrenshoog
		} else if user.YearlyIncome <= 40024 {
			return Aftoppingsgrenshoog, Maximalehuurgrens
		}
	}

	// minimum allowed to rent via woningcorporation
	return Aftoppingsgrenshoog, math.MaxInt32
}
