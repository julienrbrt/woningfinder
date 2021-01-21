package services

import "errors"

// ErrNoMatchFound is returned when there is no match found (credentials, offers,...)
var ErrNoMatchFound = errors.New("no match found")
