package matcher

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/database"
	"github.com/woningfinder/woningfinder/internal/services"
)

// MatcherOffer matcher a corporation offer with customer housing preferences
func (s *service) MatchOffer(ctx context.Context, offers corporation.Offers) error {
	// create housing corporation client
	client, err := s.clientProvider.Get(offers.Corporation.Name)
	if err != nil {
		return fmt.Errorf("error while getting corporation client %s: %w", offers.Corporation.Name, err)

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

	// match offers for each user having corporation credentials
	var wg sync.WaitGroup
	for _, creds := range credentials {
		user := &customer.User{ID: creds.UserID}

		wg.Add(1)
		// react concurrently
		go func(user *customer.User, creds customer.CorporationCredentials, wg *sync.WaitGroup) {
			defer wg.Done()

			// use one housing corporation client per user
			client := client

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
			newCreds, err := s.userService.DecryptCredentials(&creds)
			if err != nil {
				s.logger.Sugar().Errorf("error while decrypting credentials for %s: %w", user.Email, err)
				return
			}

			// login to housing corporation
			if err := client.Login(newCreds.Login, newCreds.Password); err != nil {
				if !errors.Is(err, connector.ErrAuthFailed) {
					s.logger.Sugar().Errorf("failed to login to corporation %s for %s: %w", offers.Corporation.Name, user.Email, err)
					return
				}

				// user has failed login
				s.logger.Sugar().Debugf("failed to login to corporation %s for %s: %w", offers.Corporation.Name, user.Email, err)
				if err := s.hasFailedLogin(user, newCreds); err != nil {
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

				if s.matcher.MatchOffer(*user, offer) {
					// react to offer
					if err := client.React(offer); err != nil {
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
		}(user, creds, &wg)
	}

	// wait for all match
	wg.Wait()
	return nil
}

// hasFailedLogin increased the failed login count of a corporation
// after 3 failure the login credentials of that user are deleted and the user get notified
func (s *service) hasFailedLogin(user *customer.User, credentials *customer.CorporationCredentials) error {
	// update failure count
	failureCount := credentials.FailureCount + 1

	if failureCount > 3 {
		if err := s.userService.DeleteCorporationCredentials(credentials.UserID, credentials.CorporationName); err != nil {
			return fmt.Errorf("failed to delete %s corporation credentials of user %s: %w", user.Email, credentials.CorporationName, err)
		}

		if err := s.notificationService.SendCorporationCredentialsError(user, credentials.CorporationName); err != nil {
			return fmt.Errorf("failed notifying user %s about %s corporation credentials deletion: %w", user.Email, credentials.CorporationName, err)
		}
	}

	// update failure count
	if err := s.userService.UpdateCorporationCredentialsFailureCount(credentials.UserID, credentials.CorporationName, failureCount); err != nil {
		return fmt.Errorf("failed to updating user %s %s corporation credentials login failure count: %w", user.Email, credentials.CorporationName, err)
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

func buildReactionUUID(user *customer.User, offer corporation.Offer) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + offer.Housing.Address + offer.SelectionDate.String()))
}
