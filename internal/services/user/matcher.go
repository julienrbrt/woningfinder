package user

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/domain/entity"
)

var errNoMatchFound = fmt.Errorf("no match found")

func (s *service) MatchOffer(offerList entity.OfferList) error {
	// create housing corporation client
	client, err := s.clientProvider.Get(offerList.Corporation)
	if err != nil {
		return err
	}

	// find users corporation credentials for this offers
	credentials, err := s.findMatchCredentials(offerList)
	if err != nil {
		// no users found, exit silently
		if errors.Is(err, errNoMatchFound) {
			return nil
		}
		return fmt.Errorf("error while matching offer: %w", err)
	}

	// find users with housing preferences matching offer
	users, err := s.listUsersFromCredentials(credentials)
	if err != nil {
		return fmt.Errorf("error while matching offer: %w", err)
	}

	// build credentials map
	var credentialsMap = make(map[int]*entity.CorporationCredentials, len(credentials))
	for _, c := range credentials {
		credentialsMap[c.UserID] = &c
	}

	for _, user := range users {
		// react concurrently
		go func(user entity.User) {
			//get housing preferences
			housingPreferences, err := s.GetHousingPreferences(&user)
			if err != nil {
				s.logger.Sugar().Errorf("error while getting housing preferences for user %s: %w", user.Email, err)
				return
			}
			user.HousingPreferences = *housingPreferences

			// decrypt housing corporation credentials
			creds := credentialsMap[int(user.ID)]
			creds, err = s.decryptCredentials(creds)
			if err != nil {
				s.logger.Sugar().Errorf("error while decrypting credentials for %s: %w", user.Email, err)
				return
			}

			// login to housing corporation
			if err := client.Login(creds.Login, creds.Password); err != nil {
				s.logger.Sugar().Errorf("failed to login to corporation %s for %s: %w", offerList.Corporation.Name, user.Email, err)
				return
			}

			for _, offer := range offerList.Offer {
				s.logger.Sugar().Infof("checking match of %s for %s...", offer.Housing.Address, user.Email)

				// check if we already check this offer
				uuid := buildReactionUUID(&user, offer)
				if s.hasReacted(uuid) {
					s.logger.Sugar().Debug("has already been checked... skipping.")
					continue
				}

				if user.MatchPreferences(offer) && user.MatchCriteria(offer) {
					// apply
					if err := client.ReactToOffer(offer); err != nil {
						s.logger.Sugar().Errorf("failed to react to %s with user %s: %w", offer.Housing.Address, user.Email, err)
						continue
					}

					// TODO add to queue to send mail
					s.logger.Sugar().Infof("🎉🎉🎉 WoningFinder has successfully reacted to %s on behalf of %s. 🎉🎉🎉\n", offer.Housing.Address, user.Email)
				}

				// save that we've checked the offer for the user
				s.storeReaction(uuid)
			}
		}(user)
	}

	return nil
}

// hasReacted check if a user already reacted to an offer
func (s *service) hasReacted(uuid string) bool {
	_, err := s.redisClient.Get(uuid)
	if err != nil {
		if !errors.Is(err, database.ErrRedisKeyNotFound) {
			s.logger.Sugar().Errorf("error when getting reaction: %w", err)
		}
		// does not have reacted
		return false
	}

	return true
}

// storeReaction saves that an user reacted to an offer
func (s *service) storeReaction(uuid string) {
	if err := s.redisClient.Set(uuid, true); err != nil {
		s.logger.Sugar().Errorf("error when saving reaction to redis: %w", err)
	}
}

func (s *service) findMatchCredentials(offerList entity.OfferList) ([]entity.CorporationCredentials, error) {
	var credentials []entity.CorporationCredentials
	var query = entity.CorporationCredentials{
		CorporationName: offerList.Corporation.Name,
		CorporationURL:  offerList.Corporation.URL,
	}
	if err := s.dbClient.Conn().Model(&entity.CorporationCredentials{}).Where(query).Find(&credentials).Error; err != nil {
		return nil, fmt.Errorf("error of matchOffer while getting user credentials: %w", err)
	}

	// no users found, exit silently
	if len(credentials) == 0 {
		return nil, errNoMatchFound
	}

	return credentials, nil
}

func (s *service) listUsersFromCredentials(credentials []entity.CorporationCredentials) ([]entity.User, error) {
	// get users id
	var usersID []int
	for _, c := range credentials {
		usersID = append(usersID, c.UserID)
	}

	// query each user
	var users []entity.User
	if err := s.dbClient.Conn().Where("id IN ?", usersID).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("error while getting users: %w", err)
	}

	return users, nil
}

func buildReactionUUID(user *entity.User, offer entity.Offer) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + offer.Housing.Address + offer.SelectionDate.String()))
}
