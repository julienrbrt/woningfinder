package matcher

import (
	"encoding/json"
	"fmt"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"go.uber.org/zap"
)

// SendOffers sends an offer of a housing corporation to the redis queue
func (s *service) SendOffers(offers corporation.Offers) error {
	result, err := json.Marshal(offers)
	if err != nil {
		return fmt.Errorf("error while marshaling %s offers: %w", offers.CorporationName, err)
	}

	if err := s.redisClient.Push(offersQueue, result); err != nil {
		return fmt.Errorf("error pushing %s offers to queue: %w", offers.CorporationName, err)
	}

	if err := s.verifyCorporationCities(offers); err != nil {
		return fmt.Errorf("error verifying %s cities: %w", offers.CorporationName, err)
	}

	return nil
}

// gets and verify if all cities from the offers are present the supported cities by the corporation
func (s *service) verifyCorporationCities(offers corporation.Offers) error {
	corp, err := s.connectorProvider.GetCorporation(offers.CorporationName)
	if err != nil {
		return err
	}

	cities := make(map[string]city.City)

	// get cities from offers
	for _, offer := range offers.Offer {
		// merge city names
		city := (&city.City{Name: offer.Housing.CityName}).Merge()
		cities[city.Name] = city
	}

	// check against cities from housing corporation
	for _, city := range corp.Cities {
		delete(cities, city.Name)
	}

	// no cities to add
	if len(cities) == 0 {
		return nil
	}

	// transform map to array
	var notFound []city.City
	for _, city := range cities {
		notFound = append(notFound, city)
	}

	// add cities and cities relation
	if err := s.corporationService.LinkCities(notFound, false, corp); err != nil {
		return fmt.Errorf("failing adding cities to corporation: %w", err)
	}

	s.logger.Warn(fmt.Sprintf("new cities added in %s corporation", corp.Name), zap.Int("count", len(notFound)), zap.Any("cities", notFound), zap.String("corporation", corp.Name))

	return nil
}

func (s *service) RetrieveOffers(ch chan<- corporation.Offers) error {
	for {
		offers, err := s.redisClient.BPop(offersQueue)
		if err != nil {
			return err
		}

		// Consume offers from queue
		for _, o := range offers {
			var offers corporation.Offers
			err := json.Unmarshal([]byte(o), &offers)
			if err != nil {
				s.logger.Error("error while unmarshaling offers", zap.Error(err))
				continue
			}

			// send offers to channel
			ch <- offers
		}
	}
}
