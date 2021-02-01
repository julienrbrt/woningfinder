package matcher

import (
	"math"
	"time"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// last updated in 2021

// https://www.woningmarktbeleid.nl/onderwerpen/huurtoeslag
// https://www.woningmarktbeleid.nl/actueel/nieuws/2020/11/18/inkomens--en-huurgrenzen-huurtoeslag-2021-bekend
const (
	aow   = 67
	young = 23

	kwaliteitskortingsgrens = 442.46
	aftoppingsgrenslaag     = 633.25
	aftoppingsgrenshoog     = 678.66
	maximalehuurgrens       = 752.33
)

// PassendToewijzen determines the rent range to which an user can react
// https://www.domijn.nl/over-ons/actueel/nieuws/huurtoeslag-en-passend-toewijzen-in-2021/
// https://www.dewoonplaats.nl/over-ons/actueel/201214-nieuwe-inkomensgrenzen/
// TODO add AOW logic
func PassendToewijzen(u *entity.User) (float64, float64) {
	age := time.Now().Year() - u.BirthYear

	// bypass passend toewijzen
	if u.YearlyIncome == -1 {
		return 0, math.MaxInt16
	}

	// if too rich free sector only allowed
	if u.YearlyIncome > 44655 {
		return maximalehuurgrens, math.MaxInt16
	}

	// housing allowance for people younger than 23 years old and no children
	if age < young && !u.HasChildrenSameHousing {
		return 0, kwaliteitskortingsgrens
	}

	switch {
	case u.FamilySize <= 1:
		if u.YearlyIncome <= 23725 {
			return 0, aftoppingsgrenslaag
		} else if u.YearlyIncome <= 32200 {
			return 0, aftoppingsgrenshoog
		} else if u.YearlyIncome <= 40024 {
			return aftoppingsgrenshoog, maximalehuurgrens
		}
	case u.FamilySize == 2:
		if u.YearlyIncome <= 32200 {
			return 0, aftoppingsgrenslaag
		} else if u.YearlyIncome <= 40024 {
			return aftoppingsgrenshoog, maximalehuurgrens
		}
	case u.FamilySize >= 3:
		if u.YearlyIncome <= 32200 {
			return 0, aftoppingsgrenshoog
		} else if u.YearlyIncome <= 40024 {
			return aftoppingsgrenshoog, maximalehuurgrens
		}
	}

	// minimum allowed to rent via woningcorporation
	return aftoppingsgrenshoog, math.MaxInt16
}

// MatchCriteria verifies that an user match the offer criterias
func MatchCriteria(u *entity.User, offer entity.Offer) bool {
	age := time.Now().Year() - u.BirthYear

	// checks if offer age is set and check boundaries
	if offer.MinAge > 0 && ((age < offer.MinAge) || (offer.MaxAge != 0 && age > offer.MaxAge)) {
		return false
	}

	// checks if offer family size is set and check boundaries
	if offer.MinFamilySize > 0 && (u.FamilySize < offer.MinFamilySize) || (offer.MaxFamilySize > 0 && u.FamilySize > offer.MaxFamilySize) {
		return false
	}

	// checks if offer incomes is set and check boundaries
	min, max := PassendToewijzen(u)
	if offer.Housing.Price < min || offer.Housing.Price > max {
		return false
	}

	return true
}
