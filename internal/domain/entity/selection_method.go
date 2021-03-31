package entity

// SelectionMethod defines the selection method used for a Housing Corporation in an Offer
// There is 3 supported Method: SelectionRandom, SelectionFirstComeFirstServed, SelectionRegistrationDate
type SelectionMethod string

const (
	// SelectionRandom selects a candidate from an offer randomly
	SelectionRandom SelectionMethod = "SELECTION_RANDOM"
	// SelectionFirstComeFirstServed selects first candidate that reacted to an offer
	SelectionFirstComeFirstServed SelectionMethod = "SELECTION_FIRST_COME_FIRST_SERVED"
	// SelectionRegistrationDate selects the candidate that registered the first in the housing corporation in the offer drawing
	SelectionRegistrationDate SelectionMethod = "SELECTION_REGISTRATION_DATE"
)
