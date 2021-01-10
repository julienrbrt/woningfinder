package entity

const (
	// PlanZeker is the normal plan
	PlanZeker Plan = "ZEKER"
	// PlanSneller is the high-end plan
	PlanSneller = "SNELLER"
)

// Plan defines the different plans
type Plan string

// AllowMultipleHousingPreferences checks if the plan allows multiple housing preferences
func (p *Plan) AllowMultipleHousingPreferences() bool {
	return string(*p) == PlanSneller
}

// Exists check if the plan exists
func (p *Plan) Exists() bool {
	return *p != PlanZeker && *p != PlanSneller
}
