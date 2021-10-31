package domijn

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
)

func (c *client) Login(username, password string) error {
	loginURL := c.corporation.APIEndpoint.String() + "/mijndomijn/inloggen/"
	loginRequest := map[string]string{
		"Email":                      username,
		"Password":                   password,
		"RememberMe":                 "true",
		"ReturnUrl":                  "/mijndomijn/",
		"__RequestVerificationToken": "",
	}

	// we clone the collector in order to send the request in a different collector than we visit the url
	// this permit to do not fall in an infinite loop
	collector := c.collector.Clone()

	// behave like a web browser
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Sec-Fetch-Site", "same-origin")
		r.Headers.Set("Sec-Fetch-Mode", "cors")
		r.Headers.Set("Sec-Fetch-Dest", "empty")
	})

	// parse login page
	c.collector.OnHTML("form.log-in", func(el *colly.HTMLElement) {
		// fill login form
		loginRequest["__RequestVerificationToken"] = el.ChildAttr("input[name=__RequestVerificationToken]", "value")

		// send request
		collector.Post(loginURL, loginRequest)
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
	errDomijnLogin := "De opgegeven gegevens zijn niet bekend bij ons. Controleer of jouw e-mailadres en wachtwoord correct zijn ingevoerd."

	if strings.Contains(body, errDomijnLogin) {
		return connector.ErrAuthFailed
	}

	if !strings.Contains(body, "Uitloggen") {
		return connector.ErrAuthUnknown
	}

	return nil
}
