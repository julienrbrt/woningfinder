package scheduler

import "time"

// CreateSelectionTime is a helper function that build a selection time
func CreateSelectionTime(hour, minute int) time.Time {
	return time.Date(0, 0, 0, hour, minute, 0, 0, time.UTC)
}
