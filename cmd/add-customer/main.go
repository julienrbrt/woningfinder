package main

import (
	"log"

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

	corporationService := corporation.NewService(bootstrap.DB, nil)
	userService := user.NewService(bootstrap.DB, corporationService)

	// define hardcoded user preferences
	user := &user.User{
		FullName:  "Julien Robert",
		Email:     "julien@rbrt.fr",
		BirthYear: 1998,
		HousingPreferences: user.HousingPreferences{
			Type: []corporation.HousingType{
				corporation.HousingType{
					Type: corporation.House,
				},
				corporation.HousingType{
					Type: corporation.Appartement,
				},
			},
			MinimumPrice: 400,
			MaximumPrice: 960,
			City: []corporation.City{
				{
					Name:   "Enschede",
					Region: "Overijssel",
					District: []corporation.District{
						{
							Name: "roombeek",
						},
						{
							Name: "boddenkamp",
						},
						{
							Name: "lasonder",
						},
					}},
				{Name: "Hengelo", Region: "Overijssel"},
			},
			NumberBedroom: 2,
		},
	}

	user, err = userService.Create(user)
	if err != nil {
		log.Fatal(err)
	}
}
