package templates

import (
	"fmt"
	"time"

	"github.com/matcornic/hermes/v2"
)

// WoningFinderInfo contains the email information about WoningFinder
var WoningFinderInfo = hermes.Hermes{
	Theme: new(theme),
	Product: hermes.Product{
		Name:        "Team WoningFinder",
		Link:        "https://woningfinder.nl/",
		Logo:        "https://woningfinder.nl/logo.png",
		Copyright:   fmt.Sprintf("© %d WoningFinder", time.Now().Year()),
		TroubleText: "Als de \"{ACTION}\" knop niet voor jou werkt, je kan de URL hieronder klikken of kopiëren.",
	},
}
