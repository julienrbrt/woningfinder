package matcher

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/entity"
	"github.com/woningfinder/woningfinder/internal/matcher"
	"github.com/woningfinder/woningfinder/internal/services"
)

func (s *service) MatchOffer(ctx context.Context, offers entity.OfferList) error {
	// create housing corporation client
	client, err := s.clientProvider.GetByName(offers.Corporation)
	if err != nil {
		return err
	}

	// find users corporation credentials for this offers
	credentials, err := s.userService.GetAllCorporationCredentials(offers.Corporation.Name)
	if err != nil {
		// no users found, exit silently
		if errors.Is(err, services.ErrNoMatchFound) {
			return nil
		}
		return fmt.Errorf("error while matching offer: %w", err)
	}

	// match offers
	for _, cred := range credentials {
		user := &entity.User{ID: cred.UserID}

		// react concurrently
		go func(user *entity.User, cred entity.CorporationCredentials) {
			//enrich user
			user, err = s.userService.GetUser(user)
			if err != nil {
				s.logger.Sugar().Errorf("error while getting user %s: %w", user.Email, err)
				return
			}

			// skip non paid user
			if !user.HasPaid() {
				return
			}

			// decrypt housing corporation credentials
			creds, err := s.userService.DecryptCredentials(&cred)
			if err != nil {
				s.logger.Sugar().Errorf("error while decrypting credentials for %s: %w", user.Email, err)
				return
			}

			// login to housing corporation
			if err := client.Login(creds.Login, creds.Password); err != nil {
				s.logger.Sugar().Errorf("failed to login to corporation %s for %s: %w", offers.Corporation.Name, user.Email, err)

				if err = s.hasFailedLogin(user, creds); err != nil {
					s.logger.Sugar().Warn(err)
				}

				return
			}

			for _, offer := range offers.Offer {
				s.logger.Sugar().Debugf("checking match of %s for %s...", offer.Housing.Address, user.Email)

				// check if we already check this offer
				uuid := buildReactionUUID(user, offer)
				if s.hasReacted(uuid) {
					continue
				}

				if matcher.MatchPreferences(user, offer) && matcher.MatchCriteria(user, offer) {
					// react to offer
					if err := client.ReactToOffer(offer); err != nil {
						s.logger.Sugar().Errorf("failed to react to %s with user %s: %w", offer.Housing.Address, user.Email, err)
						continue
					}

					// save match to database
					if err := s.userService.CreateHousingPreferencesMatch(user, offer, creds.CorporationName); err != nil {
						s.logger.Sugar().Errorf("failed to add %s match to user %s: %w", offer.Housing.Address, user.Email, err)
					}

					s.logger.Sugar().Infof("🎉🎉🎉 WoningFinder has successfully reacted to %s on behalf of %s. 🎉🎉🎉\n", offer.Housing.Address, user.Email)
				}

				// save that we've checked the offer for the user
				s.storeReaction(uuid)
			}
		}(user, cred)
	}

	return nil
}

// hasFailedLogin increased the failed login count of a corporation
// after 3 failure the login credentials of that user are deleted and the user get notified
func (s *service) hasFailedLogin(user *entity.User, credentials *entity.CorporationCredentials) error {
	// update failure count
	credentials.FailureCount += 1

	if credentials.FailureCount > 3 {
		if err := s.userService.DeleteCorporationCredentials(credentials.UserID, credentials.CorporationName); err != nil {
			return fmt.Errorf("failed to delete %s corporation credentials of user %s: %w", user.Email, credentials.CorporationName, err)
		}

		if err := s.notificationsService.SendCorporationCredentialsError(user, credentials.CorporationName); err != nil {
			return fmt.Errorf("failed notifying user %s about %s corporation credentials deletion: %w", user.Email, credentials.CorporationName, err)
		}
	}

	// update failure count
	if err := s.userService.CreateCorporationCredentials(credentials.UserID, *credentials); err != nil {
		return fmt.Errorf("failed to updating %s corporation credentials login failure count of user %s: %w", user.Email, credentials.CorporationName, err)
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

func buildReactionUUID(user *entity.User, offer entity.Offer) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + offer.Housing.Address + offer.SelectionDate.String()))
}
