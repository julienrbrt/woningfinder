package corporation

import (
	"fmt"

	"github.com/woningfinder/woningfinder/internal/domain/entity"
	"gorm.io/gorm/clause"
)

func (s *service) CreateOrUpdateCorporation(corporations *[]entity.Corporation) (*[]entity.Corporation, error) {
	// creates the corporation - on data changes update it
	if err := s.dbClient.Conn().Clauses(clause.OnConflict{UpdateAll: true}).Create(corporations).Error; err != nil {
		return nil, err
	}

	return corporations, nil
}

func (s *service) GetCity(name string) (*entity.City, error) {
	var c entity.City
	if err := s.dbClient.Conn().Where(entity.City{Name: name}).First(&c).Error; err != nil {
		return nil, fmt.Errorf("failing getting city %s: %w", name, err)
	}

	if c.Name == "" {
		return nil, fmt.Errorf("no city found with the name: %s", name)
	}

	return &c, nil
}
