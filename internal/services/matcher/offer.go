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

	if err := s.redisClient.Publish(pubSubOffers, result); err != nil {
		return fmt.Errorf("error publishing %d offers: %w", len(offers), err)
	}

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
