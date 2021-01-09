package user

const (
	// Zeker is the normal plan
	Zeker Plan = "ZEKER"
	// Sneller is the high-end plan
	Sneller = "SNELLER"
)

// Plan defines the different plans
type Plan string

func (u *User) canHaveMultiplePreferences() bool {
	return u.Plan == Sneller
}
