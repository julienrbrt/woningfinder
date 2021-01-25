package templates

import (
	"fmt"
	"time"

	"github.com/matcornic/hermes/v2"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

// Templates permits to build emails for users
type Templates interface {
	WelcomeTpl() (html, plain string, err error)
	WeeklyUpdateTpl(housingMatch []entity.HousingPreferencesMatch) (html, plain string, err error)
	CorporationCredentialsErrorTpl(corporationName string) (html, plain string, err error)
	ByeTpl() (html, plain string, err error)
}

type templates struct {
	product  hermes.Hermes
	user     entity.User
	jwtToken string
}

// NewTemplates builds email for WoningFinder
func NewTemplates(user entity.User, jwtToken string) Templates {
	return &templates{product: hermes.Hermes{
		Theme: new(woningFinderTheme),
		Product: hermes.Product{
			Name:        "Team WoningFinder",
			Link:        "https://woningfinder.nl/",
			Logo:        "https://cdn.umso.co/66g5xhaavhpm/assets/i0k3xo8x.png",
			Copyright:   fmt.Sprintf("© %d WoningFinder", time.Now().Year()),
			TroubleText: "Als de \"{ACTION}\" knop niet voor jou werkt, je kan de URL hieronder klikken of kopiëren.",
		},
	},
		user:     user,
		jwtToken: jwtToken,
	}
}
