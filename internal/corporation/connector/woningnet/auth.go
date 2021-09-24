package woningnet

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
)

func (c *client) Login(username, password string) error {
	loginURL := c.corporation.URL + "/Inloggen"
	loginRequest := map[string][]byte{
		"OnthoudGebruikersnaam":      []byte("false"),
		"ReturnUrl":                  nil,
		"__RequestVerificationToken": nil,
		"gebruikersnaam":             []byte(username),
		"password":                   []byte(password),
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
	c.collector.OnHTML("form", func(el *colly.HTMLElement) {
		// fill login form
		loginRequest["__RequestVerificationToken"] = []byte(el.ChildAttr("input[name=__RequestVerificationToken]", "value"))

		// send request
		collector.PostMultipart(loginURL, loginRequest)
	})

	// parse login error (from second collector)
	var hasErrLogin error
	collector.OnScraped(func(resp *colly.Response) {
		hasErrLogin = checkLogin(string(resp.Body))
	})

	// visit login page
	if err := c.collector.Visit(loginURL); err != nil {
		return fmt.Errorf("woningnet connector: error visiting login page: %w", err)
	}

	return hasErrLogin
}

func checkLogin(body string) error {
	if strings.Contains(body, "Je inloggegevens zijn niet correct.") {
		return connector.ErrAuthFailed
	}

	if !strings.Contains(body, "Uitloggen") {
		return connector.ErrAuthUnknown
	}

	return nil
}
