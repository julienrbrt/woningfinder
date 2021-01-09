package corporation

import (
	"encoding/json"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/database"
	"go.uber.org/zap"

	"gorm.io/gorm/clause"
)

// Service permits to handle the persistence of a corporation
type Service interface {
	// Create Corporation
	CreateOrUpdate(corporation *[]Corporation) (*[]Corporation, error)
	CreateHousingType(housingTypes *[]HousingType) (*[]HousingType, error)

	// Pub-Sub Offers
	PublishOffers(client Client, corporation Corporation) error
	SubscribeOffers(offerCh chan<- OfferList)

	// Getters
	GetCity(name string) (*City, error)
}

type service struct {
	logger      *zap.Logger
	dbClient    database.DBClient
	redisClient database.RedisClient
}

func NewService(logger *zap.Logger, dbClient database.DBClient, redisClient database.RedisClient) Service {
	return &service{
		logger:      logger,
		dbClient:    dbClient,
		redisClient: redisClient,
	}
}

func (s *service) CreateOrUpdate(corporations *[]Corporation) (*[]Corporation, error) {
	// creates the corporation - on data changes update it
	if err := s.dbClient.Conn().Clauses(clause.OnConflict{UpdateAll: true}).Create(corporations).Error; err != nil {
		return nil, err
	}

	return corporations, nil
}

func (s *service) CreateHousingType(housingTypes *[]HousingType) (*[]HousingType, error) {
	// creates housing types
	if err := s.dbClient.Conn().Clauses(clause.OnConflict{UpdateAll: true}).Create(housingTypes).Error; err != nil {
		return nil, err
	}

	return housingTypes, nil
}

func (s *service) GetCity(name string) (*City, error) {
	var c City
	if err := s.dbClient.Conn().Where(City{Name: name}).First(&c).Error; err != nil {
		return nil, fmt.Errorf("failing getting city %s: %w", name, err)
	}

	if c.Name == "" {
		return nil, fmt.Errorf("no city found with the name: %s", name)
	}

	return &c, nil
}

func (s *service) PublishOffers(client Client, corporation Corporation) error {
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

	if err := s.redisClient.Publish(database.PubSubOffers, result); err != nil {
		return fmt.Errorf("error publishing %d offers: %w", len(offers), err)
	}

	return nil
}

func (s *service) SubscribeOffers(offerCh chan<- OfferList) {
	ch, err := s.redisClient.Subscribe(database.PubSubOffers)
	if err != nil {
		s.logger.Sugar().Error(err)
	}

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
