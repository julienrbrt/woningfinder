package scheduler

import "time"

// CreateSelectionTime is a helper function that build a selection time
func CreateSelectionTime(hour, minute, second int) time.Time {
	return time.Date(2021, 1, 1, hour, minute, second, 0, time.UTC)
}
