package matcher

import (
	"encoding/json"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/city"
	"github.com/woningfinder/woningfinder/internal/city/manipulate"
	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
)

// PushOffers pushes the offers of a housing corporation to redis queue
func (s *service) PushOffers(client connector.Client, corp corporation.Corporation) error {
	offers, err := client.GetOffers()
	if err != nil {
		return fmt.Errorf("error while fetching offers for %s: %w", corp.Name, err)
	}

	// log number of offers found
	if len(offers) == 0 {
		s.logger.Sugar().Infof("no offers found for %s", corp.Name)
		return nil
	}

	s.logger.Sugar().Infof("%d offers found for %s", len(offers), corp.Name)

	// build offers list
	result, err := json.Marshal(corporation.Offers{
		Corporation: corp,
		Offer:       offers,
	})
	if err != nil {
		return fmt.Errorf("error while marshaling offers for %s: %w", corp.Name, err)
	}

	// add to redis queue
	if err := s.redisClient.Push(offersQueue, result); err != nil {
		return fmt.Errorf("error pushing %d offers to queue: %w", len(offers), err)
	}

	// verify supported cities
	if err := s.verifyCorporationCities(offers, corp); err != nil {
		return fmt.Errorf("error verifying corporation %s cities: %w", corp.Name, err)
	}

	return nil
}

// gets and verify if all cities from the offers are present the supported cities by the corporation
func (s *service) verifyCorporationCities(offers []corporation.Offer, corp corporation.Corporation) error {
	cities := make(map[string]city.City)
	// get cities from offers
	for _, offer := range offers {
		// merge city names
		city := manipulate.MergeCity(offer.Housing.City)
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
	if err := s.corporationService.LinkCities(notFound, corp); err != nil {
		return fmt.Errorf("failing adding cities to corporation: %w", err)
	}

	s.logger.Sugar().Warnf("%d cities (%+v) added for %s: please update housing corporation informations", len(notFound), notFound, corp.Name)

	return nil
}

func (s *service) SubscribeOffers(ch chan<- corporation.Offers) error {
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
				s.logger.Sugar().Errorf("error while unmarshaling offers: %w", err)
				continue
			}

			// send offers to channel
			ch <- offers
		}
	}
}
