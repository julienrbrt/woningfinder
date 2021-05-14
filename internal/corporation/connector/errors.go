package connector

import "errors"

var ErrAuthFailed = errors.New("authentication failed")
var ErrAuthUnknown = errors.New("error authentication: unknown error")
var ErrReactUnknown = errors.New("error reacting to house")
