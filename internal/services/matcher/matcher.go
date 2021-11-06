package matcher

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/woningfinder/woningfinder/internal/corporation"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
	"github.com/woningfinder/woningfinder/internal/customer"
	"github.com/woningfinder/woningfinder/internal/database"
	"go.uber.org/zap"
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

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)

		// skip user with invalid plan
		if !user.Plan.IsValid() {
			continue
		}

		// use one housing corporation client per user
		client := client // https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable
		go s.matchOffers(&wg, client, user, offers)
	}

	wg.Wait()
	return nil
}

// matchOffers for each user having corporation credentials
func (s *service) matchOffers(wg *sync.WaitGroup, client connector.Client, user *customer.User, offers corporation.Offers) {
	defer wg.Done()

	// check if we already checked all offers
	// this is done before login in order to do not spam login to the housing corporation and reacting to nothing
	uncheckedOffers, ok := s.hasNonReactedOffers(user, offers)
	if !ok {
		s.logger.Info("no new offers", zap.String("corporation", offers.Corporation.Name), zap.String("email", user.Email))
		return
	}

	// enrich housing preferences
	var err error
	user.HousingPreferences, err = s.userService.GetHousingPreferences(user.ID)
	if err != nil {
		s.logger.Error("error getting user preferences", zap.String("email", user.Email), zap.Error(err))
		return
	}

	// decrypt credentials
	newCreds, err := s.userService.DecryptCredentials(user.CorporationCredentials[0])
	if err != nil {
		s.logger.Error("error while decrypting credentials", zap.String("email", user.Email), zap.Error(err))
		return
	}

	// login to housing corporation
	if err := client.Login(newCreds.Login, newCreds.Password); err != nil {
		if !errors.Is(err, connector.ErrAuthFailed) {
			s.logger.Error("failed to login to corporation", zap.String("corporation", offers.Corporation.Name), zap.String("email", user.Email), zap.Error(err))
			return
		}

		// user has failed login
		s.logger.Info("failed to login to corporation", zap.String("corporation", offers.Corporation.Name), zap.String("email", user.Email), zap.Error(err))
		if err := s.hasFailedLogin(user, newCreds); err != nil {
			s.logger.Error("failed to update corporation credentials", zap.Error(err))
		}

		return
	}

	for uuid, offer := range uncheckedOffers {
		s.logger.Debug("checking match", zap.String("address", offer.Housing.Address), zap.String("corporation", offers.Corporation.Name), zap.String("email", user.Email))

		if s.matcher.MatchOffer(*user, offer) {
			// react to offer
			if err := client.React(offer); err != nil {
				s.logger.Debug("failed to react", zap.String("address", offer.Housing.Address), zap.String("corporation", offers.Corporation.Name), zap.String("email", user.Email), zap.Error(err))

				// check if we retry next time or mark the offer as checked
				if ok := s.retryReactNextTime(uuid); !ok {
					s.redisClient.SetUUID(uuid)

					// alert user
					if user.HasAlertsEnabled {
						if err := s.emailService.SendReactionFailure(user, offers.Corporation.Name, offer); err != nil {
							s.logger.Error("failed to send email", zap.Error(err))
						}
					}
				}

				continue
			}

			// get and upload housing picture
			pictureURL := s.uploadHousingPicture(offer)

			// save match to database
			if err := s.userService.CreateHousingPreferencesMatch(user.ID, offer, user.CorporationCredentials[0].CorporationName, pictureURL); err != nil {
				s.logger.Error("failed to add housing preferences match", zap.String("address", offer.Housing.Address), zap.String("corporation", offers.Corporation.Name), zap.String("email", user.Email), zap.Error(err))
			}

			s.logger.Info("🎉🎉🎉 WoningFinder has successfully reacted to a house", zap.String("address", offer.Housing.Address), zap.String("corporation", offers.Corporation.Name), zap.String("email", user.Email), zap.Error(err))
		}

		// save that we've checked the offer for the user
		s.redisClient.SetUUID(uuid)
	}
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
		return fmt.Errorf("failed updating user %s %s corporation credentials login failure count: %w", user.Email, credentials.CorporationName, err)
	}

	return nil
}

// hasNonReactedOffers returns the offers that has not been already reacted to
func (s *service) hasNonReactedOffers(user *customer.User, offers corporation.Offers) (map[string]corporation.Offer, bool) {
	uncheckedOffers := make(map[string]corporation.Offer)

	for _, offer := range offers.Offer {
		uuid := buildReactionUUID(user, offer)
		if s.redisClient.HasUUID(uuid) {
			continue
		}

		uncheckedOffers[uuid] = offer
	}

	return uncheckedOffers, len(uncheckedOffers) > 0
}

func (s *service) uploadHousingPicture(offer corporation.Offer) string {
	fileName, err := s.spacesClient.UploadPicture("offers", offer.Housing.Address, offer.RawPictureURL)
	if err != nil {
		s.logger.Error("failed to upload picture", zap.Error(err))
	}

	return fileName
}

// retryReactNextTime checks if a offer can still be retried in a next check
// after 3 retries it returns false as the maximum of retries if 3
func (s *service) retryReactNextTime(uuid string) bool {
	maxRetry := 3
	failedUUID := "failed" + uuid

	failureCount, err := s.redisClient.Get(failedUUID)
	if err != nil {
		if !errors.Is(err, database.ErrRedisKeyNotFound) {
			s.logger.Error("error when getting reaction failure count from redis", zap.Error(err))
			return true
		}

		s.redisClient.Set(failedUUID, 1)
		return true
	}

	failureCountInt, err := strconv.Atoi(failureCount)
	if err != nil {
		s.logger.Error("error when converting reaction failure count to int", zap.Error(err))
		return true
	}

	// stop reacting to house after 3 failures
	if failureCountInt < maxRetry {
		s.redisClient.Set(failedUUID, failureCountInt+1)
		return true
	}

	return false
}

func buildReactionUUID(user *customer.User, offer corporation.Offer) string {
	return base64.StdEncoding.EncodeToString([]byte(user.Email + offer.Housing.Address))
}
