package matcher

import (
	"math"
	"time"

	"github.com/woningfinder/woningfinder/internal/customer"
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
func (m *matcher) passendToewijzen(user customer.User) (float64, float64) {
	age := time.Now().Year() - user.BirthYear

	// bypass passend toewijzen
	if user.YearlyIncome == -1 {
		return 0, math.MaxInt16
	}

	// if too rich free sector only allowed
	if user.YearlyIncome > 44655 {
		return maximalehuurgrens, math.MaxInt16
	}

	// housing allowance for people younger than 23 years old and no children
	if age < young && !user.HasChildrenSameHousing {
		return 0, kwaliteitskortingsgrens
	}

	switch {
	case user.FamilySize <= 1:
		if user.YearlyIncome <= 23725 {
			return 0, aftoppingsgrenslaag
		} else if user.YearlyIncome <= 32200 {
			return 0, aftoppingsgrenshoog
		} else if user.YearlyIncome <= 40024 {
			return aftoppingsgrenshoog, maximalehuurgrens
		}
	case user.FamilySize == 2:
		if user.YearlyIncome <= 32200 {
			return 0, aftoppingsgrenslaag
		} else if user.YearlyIncome <= 40024 {
			return aftoppingsgrenshoog, maximalehuurgrens
		}
	case user.FamilySize >= 3:
		if user.YearlyIncome <= 32200 {
			return 0, aftoppingsgrenshoog
		} else if user.YearlyIncome <= 40024 {
			return aftoppingsgrenshoog, maximalehuurgrens
		}
	}

	// minimum allowed to rent via woningcorporation
	return aftoppingsgrenshoog, math.MaxInt16
}
