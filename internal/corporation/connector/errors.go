package connector

import "errors"

var ErrAuthFailed = errors.New("authentication failed")
var ErrAuthUnknown = errors.New("unknown authentication error")
var ErrReactUnknown = errors.New("error reacting to house")
