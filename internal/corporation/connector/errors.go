package connector

import "errors"

// ErrAuthFailed is retuned when the authentication to the housing coporation has failed
var ErrAuthFailed = errors.New("authentication failed")
