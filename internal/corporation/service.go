package corporation

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"go.uber.org/zap"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const pubSubChannelName = "offers"

// Service permits to handle the persistence of a corporation
type Service interface {
	// Create Corporation
	CreateOrUpdate(corporation *[]Corporation) (*[]Corporation, error)
	CreateHousingType(housingTypes *[]HousingType) (*[]HousingType, error)

	// (Redis) Pub-Sub
	PublishOffers(client Client, corporation Corporation) error
	SubscribeOffers(offerCh chan<- OfferList)

	// Getters
	GetCity(name string) (*City, error)
}

type corporationService struct {
	logger *zap.Logger
	db     *gorm.DB
	rdb    *redis.Client
}

func NewService(logger *zap.Logger, db *gorm.DB, rdb *redis.Client) Service {
	return &corporationService{
		logger: logger,
		db:     db,
		rdb:    rdb,
	}
}

func (s *corporationService) CreateOrUpdate(corporations *[]Corporation) (*[]Corporation, error) {
	// creates the corporation - on data changes update it
	if err := s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(corporations).Error; err != nil {
		return nil, err
	}

	return corporations, nil
}

func (s *corporationService) CreateHousingType(housingTypes *[]HousingType) (*[]HousingType, error) {
	// creates housing types
	if err := s.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(housingTypes).Error; err != nil {
		return nil, err
	}

	return housingTypes, nil
}

func (s *corporationService) GetCity(name string) (*City, error) {
	var c City
	if err := s.db.Where(City{Name: name}).First(&c).Error; err != nil {
		return nil, fmt.Errorf("failing getting city %s: %w", name, err)
	}

	if c.Name == "" {
		return nil, fmt.Errorf("no city found with the name: %s", name)
	}

	return &c, nil
}

func (s *corporationService) PublishOffers(client Client, corporation Corporation) error {
	offers, err := client.FetchOffer()
	if err != nil {
		return fmt.Errorf("error while fetching offers for %s: %w", corporation.Name, err)
	}

	// log number of offers found
	if len(offers) > 0 {
		s.logger.Sugar().Infof("%d offers found for %s", len(offers), corporation.Name)
	} else {
		s.logger.Sugar().Infof("no offers found for %s", corporation.Name)
	}

	// build offers list
	offerList := OfferList{
		Corporation: corporation,
		Offer:       offers,
	}

	result, err := json.Marshal(offerList)
	if err != nil {
		return fmt.Errorf("erorr while marshaling offers for %s: %w", corporation.Name, err)
	}

	if err = s.rdb.Publish(pubSubChannelName, result).Err(); err != nil {
		return fmt.Errorf("error publishing %d offers to channel %s: %w", len(offers), pubSubChannelName, err)
	}

	return nil
}

func (s *corporationService) SubscribeOffers(offerCh chan<- OfferList) {
	pubsub := s.rdb.Subscribe(pubSubChannelName)
	defer pubsub.Close()

	// Wait for confirmation that subscription is created before doing anything.
	_, err := pubsub.Receive()
	if err != nil {
		s.logger.Sugar().Errorf("error subscribing to channel: %w", err)
	}

	// Go channel which receives messages.
	ch := pubsub.Channel()
	// Consume messages
	for msg := range ch {
		var offerList OfferList
		err := json.Unmarshal([]byte(msg.Payload), &offerList)
		if err != nil {
			s.logger.Sugar().Errorf("error while unmarshaling offers: %w", err)
			continue
		}

		go func(offers OfferList) { offerCh <- offers }(offerList)
	}
}
