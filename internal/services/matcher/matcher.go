package matcher

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/julienrbrt/woningfinder/internal/corporation"
	"github.com/julienrbrt/woningfinder/internal/corporation/city"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"github.com/julienrbrt/woningfinder/internal/customer"
	"go.uber.org/zap"
)

// Match matches corporation offers with customer housing preferences
func (s *service) Match(ctx context.Context, offers corporation.Offers) error {
	// add missing cities to corporations
	if err := s.verifyCorporationCities(offers); err != nil {
		return fmt.Errorf("error verifying %s cities: %w", offers.CorporationName, err)
	}

	// create housing corporation client
	client, err := s.connectorProvider.GetClient(offers.CorporationName)
	if err != nil {
		return fmt.Errorf("error while getting corporation client %s: %w", offers.CorporationName, err)
	}

	// find users corporation credentials for this offers
	users, err := s.userService.GetUsersWithGivenCorporationCredentials(offers.CorporationName)
	if err != nil {
		return fmt.Errorf("error while matching offer: %w", err)
	}

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)

		// skip inactive user
		if !user.IsActivated() {
			continue
		}

		// use one housing corporation client per user
		go s.matchOffers(&wg, client, user, offers)
	}

	wg.Wait()
	return nil
}

// matchOffers for each user having corporation credentials
func (s *service) matchOffers(wg *sync.WaitGroup, client connector.Client, user *customer.User, offers corporation.Offers) {
	defer wg.Done()

	// enrich housing preferences
	var err error
	user.HousingPreferences, err = s.userService.GetHousingPreferences(user.ID)
	if err != nil {
		s.logger.Error("error getting user preferences", zap.String("email", user.Email), zap.Error(err))
		return
	}

	// gets the user matching offers
	// done before any query in order to do not spam login
	matchingOffers, ok := s.getMatchingOffers(user, offers)
	if !ok {
		return
	}

	s.logger.Info(fmt.Sprintf("users matched with %d offers", len(matchingOffers)), zap.String("email", user.Email))

	// decrypt credentials
	newCreds, err := s.userService.DecryptCredentials(user.CorporationCredentials[0])
	if err != nil {
		s.logger.Error("error while decrypting credentials", zap.String("email", user.Email), zap.Error(err))
		return
	}

	// login to housing corporation
	if err := client.Login(newCreds.Login, newCreds.Password); err != nil {
		if !errors.Is(err, connector.ErrAuthFailed) {
			s.logger.Error("failed to login to corporation", zap.String("corporation", offers.CorporationName), zap.String("email", user.Email), zap.Error(err))
			return
		}

		// user has failed login
		s.logger.Info("failed to login to corporation", zap.String("corporation", offers.CorporationName), zap.String("email", user.Email), zap.Error(err))
		if err := s.updateFailedLogin(user, newCreds); err != nil {
			s.logger.Error("failed to update corporation credentials", zap.Error(err))
		}

		return
	}

	// failedReaction keeps track of the fail reaction in order to send a bulk alert
	failedReaction := []corporation.Offer{}

	for uuid, offer := range matchingOffers {
		// react to offer
		if err := client.React(offer); err != nil {
			s.logger.Info("failed to react", zap.String("url", offer.URL), zap.String("corporation", offers.CorporationName), zap.String("email", user.Email), zap.Error(err))

			// check if we retry next time or mark the offer as checked
			if ok := s.retryReactNextTime(uuid); !ok {
				failedReaction = append(failedReaction, offer)
			}

			continue
		}

		// get and upload housing picture
		pictureURL := s.uploadHousingPicture(offer)

		// save match to database
		if err := s.userService.CreateHousingPreferencesMatch(user.ID, offer, user.CorporationCredentials[0].CorporationName, pictureURL); err != nil {
			s.logger.Error("failed to add housing preferences match", zap.String("address", offer.Housing.Address), zap.String("corporation", offers.CorporationName), zap.String("email", user.Email), zap.Error(err))
		}

		// mark the offer as checked
		s.setMatch(&MatcherCounter{UUID: uuid, Reacted: true})

		s.logger.Info("🎉🎉🎉 WoningFinder has successfully reacted to a house", zap.String("address", offer.Housing.Address), zap.String("corporation", offers.CorporationName), zap.String("email", user.Email), zap.Error(err))
	}

	// alert user in case of failed reaction
	if user.HasAlertsEnabled && len(failedReaction) > 0 {
		if err := s.emailService.SendReactionFailure(user, offers.CorporationName, failedReaction); err != nil {
			s.logger.Error("failed to send reaction failure email", zap.Error(err))
		}
	}
}

// updateFailedLogin increased the failed login count of a corporation
// after 3 failure the login credentials of that user are deleted and the user get notified
func (s *service) updateFailedLogin(user *customer.User, credentials *customer.CorporationCredentials) error {
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

// getMatchingOffers returns the user matching offers
func (s *service) getMatchingOffers(user *customer.User, offers corporation.Offers) (map[string]corporation.Offer, bool) {
	matchingOffers := make(map[string]corporation.Offer)

	for _, offer := range offers.Offer {
		uuid := buildReactionUUID(user, offer)
		if !s.matcher.MatchOffer(*user, offer) {
			continue
		}

		if _, err := s.getMatch(uuid); err != nil {
			continue
		}

		matchingOffers[uuid] = offer
	}

	return matchingOffers, len(matchingOffers) > 0
}

func (s *service) uploadHousingPicture(offer corporation.Offer) string {
	fileName := strings.ToLower(hex.EncodeToString([]byte(offer.Housing.Address)))
	if err := s.imgClient.Download(fileName, offer.RawPictureURL); err != nil {
		s.logger.Error("failed to upload picture", zap.Error(err))
	}

	return fileName
}

// retryReactNextTime checks if an offer can still be retried in a next check.
// After 3 retries it returns false as the maximum of retries reaches 8
func (s *service) retryReactNextTime(uuid string) bool {
	match, err := s.getMatch(uuid)
	if err != nil {
		s.logger.Error("error when getting reaction match from database", zap.Error(err))
		return true
	}

	// stop reacting to house after 8 failures
	if match.FailureCount < 8 {
		match.FailureCount++
		s.setMatch(match)
		return true
	}

	// mark the offer as checked to not retry again
	match.Reacted = true
	s.setMatch(match)
	return false
}

func buildReactionUUID(user *customer.User, offer corporation.Offer) string {
	return hex.EncodeToString([]byte(user.Email + offer.Housing.Address))
}

// gets and verify if all cities from the offers are present the supported cities by the corporation
func (s *service) verifyCorporationCities(offers corporation.Offers) error {
	corp, err := s.connectorProvider.GetCorporation(offers.CorporationName)
	if err != nil {
		return err
	}

	cities := make(map[string]city.City)

	// get cities from offers
	for _, offer := range offers.Offer {
		// merge city names
		city := (&city.City{Name: offer.Housing.CityName}).Merge()
		cities[city.Name] = city
	}

	// check against cities from housing corporation
	for _, city := range corp.Cities {
		delete(cities, city.Name)
	}

	// no cities to add
	if len(cities) == 0 {
		return nil
	}

	// transform map to array
	var notFound []city.City
	for _, city := range cities {
		notFound = append(notFound, city)
	}

	// add cities and cities relation
	if err := s.corporationService.LinkCities(notFound, false, corp); err != nil {
		return fmt.Errorf("failing adding cities to corporation: %w", err)
	}

	s.logger.Warn(fmt.Sprintf("new cities added in %s corporation", corp.Name), zap.Int("count", len(notFound)), zap.Any("cities", notFound), zap.String("corporation", corp.Name))

	return nil
}
