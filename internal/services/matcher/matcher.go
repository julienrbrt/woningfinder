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
)

// MatcherOffer matcher a corporation offer with customer housing preferences
func (s *service) MatchOffer(ctx context.Context, offers corporation.Offers) error {
	// create housing corporation client
	client, err := s.clientProvider.Get(offers.Corporation.Name)
	if err != nil {
		return fmt.Errorf("error while getting corporation client %s: %w", offers.Corporation.Name, err)

	}

	// find users corporation credentials for this offers
	users, err := s.userService.GetUsersWithGivenCorporationCredentials(offers.Corporation.Name)
	if err != nil {
		return fmt.Errorf("error while matching offer: %w", err)
	}

	// match offers for each user having corporation credentials
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		// react concurrently
		go func(wg *sync.WaitGroup, user *customer.User) {
			defer wg.Done()

			// use one housing corporation client per user
			client := client()

			// skip user with invalid plan (not paid and free trial-expired)
			if !user.Plan.IsValid() {
				return
			}

			// enrinch housing preferences
			user.HousingPreferences, err = s.userService.GetHousingPreferences(user.ID)
			if err != nil {
				s.logger.Sugar().Warnf("failed to get users preferences from %s: %w", user.Email, err)
				return
			}

			// check if we already checked all offers
			// this is done before login in order to do not spam login to the housing corporation and reacting to nothing
			uncheckedOffers, ok := s.hasNonReactedOffers(user, offers)
			if !ok {
				s.logger.Sugar().Debugf("no new offers from %s for %s...", offers.Corporation.Name, user.Email)
				return
			}

			newCreds, err := s.userService.DecryptCredentials(user.CorporationCredentials[0])
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

			for uuid, offer := range uncheckedOffers {
				s.logger.Sugar().Debugf("checking match of %s from %s for %s...", offer.Housing.Address, offers.Corporation.Name, user.Email)

				if s.matcher.MatchOffer(*user, offer) {
					// react to offer
					if err := client.React(offer); err != nil {
						s.logger.Sugar().Errorf("failed to react to %s from %s with user %s: %w", offer.Housing.Address, offers.Corporation.Name, user.Email, err)
						continue
					}

					// get and upload housing picture
					pictureURL := s.uploadHousingPicture(offer)

					// save match to database
					if err := s.userService.CreateHousingPreferencesMatch(user.ID, offer, user.CorporationCredentials[0].CorporationName, pictureURL); err != nil {
						s.logger.Sugar().Errorf("failed to add %s from %s match to user %s: %w", offer.Housing.Address, offers.Corporation.Name, user.Email, err)
					}

					s.logger.Sugar().Infof("🎉🎉🎉 WoningFinder has successfully reacted to %s from %s on behalf of %s. 🎉🎉🎉\n", offer.Housing.Address, offers.Corporation.Name, user.Email)
				}

				// save that we've checked the offer for the user
				s.storeReaction(uuid)
			}
		}(&wg, user)
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
			return fmt.Errorf("failed to delete %s corporation credentials of user %s: %w", credentials.CorporationName, user.Email, err)
		}

		if err := s.emailService.SendCorporationCredentialsError(user, credentials.CorporationName); err != nil {
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

// hasNonReactedOffers returns the offers that has not been already reacted to
func (s *service) hasNonReactedOffers(user *customer.User, offers corporation.Offers) (map[string]corporation.Offer, bool) {
	uncheckedOffers := make(map[string]corporation.Offer)

	for _, offer := range offers.Offer {
		uuid := buildReactionUUID(user, offer)
		if s.hasReacted(uuid) {
			continue
		}

		uncheckedOffers[uuid] = offer
	}

	return uncheckedOffers, len(uncheckedOffers) > 0
}

func (s *service) uploadHousingPicture(offer corporation.Offer) string {
	fileName, err := s.spacesClient.UploadOfferPicture(offer.Housing.Address, offer.RawPictureURL)
	if err != nil {
		s.logger.Sugar().Warnf("failed to upload picture: %w", err)
	}

	return fileName
}

func buildReactionUUID(user *customer.User, offer corporation.Offer) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + offer.Housing.Address))
}
