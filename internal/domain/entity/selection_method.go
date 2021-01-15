package entity

import (
	"database/sql/driver"
)

// SelectionMethod is the database representation of Method
type SelectionMethod struct {
	Method Method `pg:",pk"`
}

const (
	// SelectionRandom selects a candidate from an offer randomly
	SelectionRandom Method = "SELECTION_RANDOM"
	// SelectionFirstComeFirstServed selects first candidate that reacted to an offer
	SelectionFirstComeFirstServed = "SELECTION_FIRST_COME_FIRST_SERVED"
	// SelectionRegistrationDate selects the candidate that registered the first in the housing corporation in the offer drawing
	SelectionRegistrationDate = "SELECTION_REGISTRATION_DATE"
)

// Method defines the selection method used for a Housing Corporation in an Offer
// There is 3 supported Method: SelectionRandom, SelectionFirstComeFirstServed, SelectionRegistrationDate
type Method string

// Scan implements the Scanner interface from reading from the database
func (m *Method) Scan(value interface{}) error {
	*m = Method(string(value.([]byte)))
	return nil
}

// Value implements the Valuer interface for the storing in the database
func (m Method) Value() (driver.Value, error) {
	return string(m), nil
}

// Exists returns is the given method exist
func (m Method) Exists() bool {
	return m == SelectionRandom || m == SelectionFirstComeFirstServed || m == SelectionRegistrationDate
}
