package retry

import "time"

// Sleeper sleeps for a given time.Duration.
type Sleeper func(time.Duration)
