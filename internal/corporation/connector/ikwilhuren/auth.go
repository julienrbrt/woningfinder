package ikwilhuren

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/julienrbrt/woningfinder/internal/corporation/connector"
	"go.uber.org/zap"
)

func (c *client) Login(username, password string) error {
	loginURL := c.corporation.URL + "/login"
	loginRequest := map[string][]byte{
		"csrf":          nil,
		"postrequest":   []byte("login"),
		"txtEmail":      []byte(username),
		"txtWachtwoord": []byte(password),
	}

	// we clone the collector in order to send the request in a different collector than we visit the url
	// this permit to do not fall in an infinite loop
	collector := c.collector.Clone()

	// behave like a web browser
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Sec-Fetch-Site", "same-origin")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
	})

	// parse login page
	c.collector.OnHTML("form", func(h *colly.HTMLElement) {
		// fill login form
		csrf := h.ChildAttr("input[name=csrf]", "value")

		loginRequest["csrf"] = []byte(csrf)

		// send request
		if err := collector.PostMultipart(loginURL, loginRequest); err != nil {
			c.logger.Error("error sending login request", logConnector, zap.Error(err))
		}
	})

	// parse login error (from second collector)
	var hasErrLogin error
	collector.OnScraped(func(resp *colly.Response) {
		hasErrLogin = c.checkLogin(string(resp.Body))
	})

	// visit login page
	if err := c.collector.Visit(loginURL); err != nil {
		return fmt.Errorf("error visiting login page: %w", err)
	}

	return hasErrLogin
}

func (c *client) checkLogin(body string) error {
	if strings.Contains(body, "De opgegeven login/wachtwoord combinatie is ongeldig.") {
		return connector.ErrAuthFailed
	}

	if !strings.Contains(body, "Uitloggen") {
		return connector.ErrAuthUnknown
	}

	return nil
}
