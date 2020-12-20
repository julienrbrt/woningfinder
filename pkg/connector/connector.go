package connector

import "github.com/woningfinder/woningfinder/internal/corporation"

// Connector specifies information a connector to a ERP
type Connector interface {
	corporation.Client
}
