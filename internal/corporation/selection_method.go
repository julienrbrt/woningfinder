package corporation

import (
	"database/sql/driver"
	"time"

	"gorm.io/gorm"
)

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
func (u *Method) Scan(value interface{}) error {
	*u = Method(value.(string))
	return nil
}

// Value implements the Valuer interface for the storing in the database
func (u Method) Value() (driver.Value, error) {
	return string(u), nil
}

// SelectionMethod is the database representation of Method
type SelectionMethod struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Method    Method         `gorm:"primaryKey"`
}
