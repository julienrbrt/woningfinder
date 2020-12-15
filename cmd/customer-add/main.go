package main

import (
	"log"
	"os"

	"github.com/woningfinder/woningfinder/internal/bootstrap"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/user"
	"github.com/woningfinder/woningfinder/pkg/env"

	"github.com/joho/godotenv"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	// fallback to system env if unexisting
	// if not defined on system, panics
	if err := godotenv.Load("../../.env"); err != nil {
		_ = env.MustGetString("APP_NAME")
	}
}

func main() {
	err := bootstrap.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	clientProvider := bootstrap.CreateClientProvider()
	corporationService := corporation.NewService(bootstrap.DB, nil)
	userService := user.NewService(bootstrap.DB, bootstrap.RDB, os.Getenv("AES_SECRET"), clientProvider, corporationService)

	// define hardcoded user preferences
	u := user.User{
		Name:         "Julien Robert",
		Email:        "julien@rbrt.fr",
		BirthYear:    1998,
		YearlyIncome: 33000,
		FamilySize:   2,
		HousingPreferences: user.HousingPreferences{
			Type: []corporation.HousingType{
				corporation.HousingType{
					Type: corporation.House,
				},
				corporation.HousingType{
					Type: corporation.Appartement,
				},
			},
			MaximumPrice: 960,
			City: []corporation.City{
				{
					Name: "Enschede",
					District: []corporation.CityDistrict{
						{Name: "Roombeek"},
						{Name: "Boddenkamp"},
						{Name: "Lasonder"},
						{Name: "Mekkelholt"},
					},
				},
				{Name: "Hengelo"},
			},
			NumberBedroom: 2,
		},
	}

	// create user
	userPtr, err := userService.CreateUser(&u)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("customer %s successfully created 🎉\n", u.Name)

	addCustomerCredentials(userService, userPtr)
}

func addCustomerCredentials(userService user.Service, u *user.User) {
	// add credentials
	creds := user.CorporationCredentials{
		Corporation: corporation.Corporation{
			Name: "De Woonplaats",
			URL:  "https://dewoonplaats.nl",
		},
		Login:    "",
		Password: "",
	}
	err := userService.CreateCorporationCredentials(u, creds)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("credentials of %s for customer %s successfully created 🎉\n", creds.Corporation.Name, u.Name)
}
