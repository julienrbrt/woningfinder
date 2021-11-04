package itris

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
)

var errItrisBlocked = errors.New("error itris authentication: woningfinder blocked")

func (c *client) Login(username, password string) error {
	loginURL := c.corporation.APIEndpoint.String() + "/inloggen/index.xml"
	loginRequest := map[string][]byte{
		"Password":                               []byte(password),
		"Username":                               []byte(username),
		"__from":                                 nil,
		"inloggen":                               []byte("inloggen"),
		"inloggen_csrf_protection":               nil,
		"inloggen_spam_protection":               []byte("PROTECT_ME"),
		"magic_roxen_automatic_charset_variable": []byte("åäö芟@UTF-8"),
		"post0":                                  nil,
		"redirect":                               nil,
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
		c.itrisCSRFToken = el.ChildAttr("input[name=inloggen_csrf_protection]", "value")
		loginRequest["inloggen_csrf_protection"] = []byte(c.itrisCSRFToken)
		loginRequest["post0"] = []byte(el.ChildAttr("input[name=post0]", "value"))

		// send request
		collector.PostMultipart(loginURL, loginRequest)
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
	errItrisLoginMsg := "Combinatie inlognaam / wachtwoord is niet bekend of onjuist. Controleer de invoer en probeer het opnieuw."
	errItrisBlockedMsg := "Om veiligheidsredenen is dit veld tijdelijk geblokkeerd, probeer het later nog eens"
	errItrisBlockedMsg2 := "De beveiliging van dit formulier weigert uw verzoek"

	if strings.Contains(body, errItrisLoginMsg) {
		return connector.ErrAuthFailed
	}

	if strings.Contains(body, errItrisBlockedMsg) || strings.Contains(body, errItrisBlockedMsg2) {
		return errItrisBlocked
	}

	if !strings.Contains(body, "Uitloggen") {
		return connector.ErrAuthUnknown
	}

	return nil
}
