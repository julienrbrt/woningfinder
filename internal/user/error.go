package user

import "fmt"

var errCorporationCredentialsNotFound = fmt.Errorf("error when getting corporation credentials: credentials not found")
var errNoMatchFound = fmt.Errorf("no match found")
