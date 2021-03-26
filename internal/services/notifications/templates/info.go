package templates

import (
	"fmt"
	"time"

	"github.com/matcornic/hermes/v2"
)

// WoningFinderInfo contains the email information about WoningFinder
var WoningFinderInfo = hermes.Hermes{
	Theme: new(woningFinderTheme),
	Product: hermes.Product{
		Name:        "Team WoningFinder",
		Link:        "https://woningfinder.nl/",
		Logo:        "https://cdn.umso.co/66g5xhaavhpm/assets/i0k3xo8x.png",
		Copyright:   fmt.Sprintf("© %d WoningFinder", time.Now().Year()),
		TroubleText: "Als de \"{ACTION}\" knop niet voor jou werkt, je kan de URL hieronder klikken of kopiëren.",
	},
}