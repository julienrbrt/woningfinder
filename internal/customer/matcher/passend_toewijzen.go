package matcher

import (
	"math"
	"time"

	"github.com/julienrbrt/woningfinder/internal/customer"
)

// last updated in november 2024

// https://www.volkshuisvestingnederland.nl/onderwerpen/huurtoeslag
// https://www.volkshuisvestingnederland.nl/onderwerpen/dossier-woningtoewijzing/documenten/publicaties/2020/12/18/infographic-toewijzen-van-woningen
const (
	aow   = 67
	young = 23

	MaximumIncomeSocialHouse = 51136
	Kwaliteitskortingsgrens  = 454.47
	Aftoppingsgrenslaag      = 650.43
	Aftoppingsgrenshoog      = 697.07
	Maximalehuurgrens        = 879.66
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
		if user.YearlyIncome <= 27725 {
			return 0, Aftoppingsgrenslaag
		} else if user.YearlyIncome <= 47699 {
			return Aftoppingsgrenslaag, Aftoppingsgrenshoog
		}
	case user.FamilySize == 2:
		if user.YearlyIncome <= 37625 {
			return 0, Aftoppingsgrenslaag
		} else if user.YearlyIncome <= 52671 {
			return Aftoppingsgrenslaag, Aftoppingsgrenshoog
		}
	case user.FamilySize >= 3:
		if user.YearlyIncome <= 37625 {
			return 0, Aftoppingsgrenshoog
		} else if user.YearlyIncome <= 52671 {
			return Aftoppingsgrenshoog, Maximalehuurgrens
		}
	}

	// minimum allowed to rent via woningcorporation
	return Aftoppingsgrenshoog, math.MaxInt32
}
