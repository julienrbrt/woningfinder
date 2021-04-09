package matcher

import (
	"encoding/json"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/entity"
)

func (s *service) PublishOffers(client corporation.Client, corporation entity.Corporation) error {
	offers, err := client.FetchOffer()
	if err != nil {
		return fmt.Errorf("error while fetching offers for %s: %w", corporation.Name, err)
	}

	// log number of offers found
	if len(offers) == 0 {
		s.logger.Sugar().Infof("no offers found for %s", corporation.Name)
		return nil
	}

	s.logger.Sugar().Infof("%d offers found for %s", len(offers), corporation.Name)

	// build offers list
	offerList := entity.OfferList{
		Corporation: corporation,
		Offer:       offers,
	}

	result, err := json.Marshal(offerList)
	if err != nil {
		return fmt.Errorf("error while marshaling offers for %s: %w", corporation.Name, err)
	}

	// add to redis queue
	if err := s.redisClient.Push(offersQueue, result); err != nil {
		return fmt.Errorf("error pushing %d offers to queue: %w", len(offers), err)
	}

	// verify supported cities
	if err := s.verifyCorporationCities(offers, corporation); err != nil {
		return fmt.Errorf("error verifying corporation %s cities: %w", corporation.Name, err)
	}

	return nil
}

// gets and verify if all cities from the offers are present the supported cities by the corporation
func (s *service) verifyCorporationCities(offers []entity.Offer, corporation entity.Corporation) error {
	cities := make(map[string]entity.City)
	// get cities from offers
	for _, offer := range offers {
		city := offer.Housing.City
		cities[city.Name] = city
	}

	// check against cities from housing corporation
	for _, city := range corporation.Cities {
		_, ok := cities[city.Name]
		if ok {
			delete(cities, city.Name)
		}
	}

	// no cities to add
	if len(cities) == 0 {
		return nil
	}

	// transform map to array
	var notFound []entity.City
	for _, city := range cities {
		notFound = append(notFound, city)
	}

	// add cities and cities relation
	if err := s.corporationService.AAACities(notFound, corporation); err != nil {
		return fmt.Errorf("failing adding cities to corporation: %w", err)
	}

	s.logger.Sugar().Warnf("%d cities (%+v) added for %s: please update housing corporation informations", len(notFound), notFound, corporation.Name)

	return nil
}

func (s *service) SubscribeOffers(ch chan<- entity.OfferList) error {
	for {
		offers, err := s.redisClient.BPop(offersQueue)
		if err != nil {
			return err
		}

		// Consume offers from queue
		for _, o := range offers {
			var offers entity.OfferList
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
