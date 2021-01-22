package matcher

import (
	"encoding/json"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
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

	// publish to redis queue
	if err := s.redisClient.Publish(pubSubOffers, result); err != nil {
		return fmt.Errorf("error publishing %d offers: %w", len(offers), err)
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
	if err := s.corporationService.AddCities(notFound, corporation); err != nil {
		return fmt.Errorf("failing adding cities to corporation: %w", err)
	}

	s.logger.Sugar().Warnf("%d cities added for %s: please update the structs", len(notFound), corporation.Name)

	return nil
}

func (s *service) SubscribeOffers(offerCh chan<- entity.OfferList) error {
	ch, err := s.redisClient.Subscribe(pubSubOffers)
	if err != nil {
		return err
	}

	// Consume messages and match offers
	for msg := range ch {
		var offers entity.OfferList
		err := json.Unmarshal([]byte(msg.Payload), &offers)
		if err != nil {
			s.logger.Sugar().Errorf("error while unmarshaling offers: %w", err)
			continue
		}

		// send offers to channel
		offerCh <- offers
	}

	// should never happen
	return nil
}
