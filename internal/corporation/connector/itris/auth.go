package itris

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/woningfinder/woningfinder/internal/corporation/connector"
)

var (
	ErrItrisLoginUnknown = errors.New("itris connector: error authentication: unknown error")
	ErrItrisBlocked      = errors.New("itris connector: error authentication: woningfinder blocked")
)

func (c *client) Login(username, password string) error {
	collector := c.collector
	collector.AllowURLRevisit = false

	loginURL := c.url + "/inloggen/index.xml"
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

	// parse login page
	collector.OnHTML("form", func(el *colly.HTMLElement) {
		// fill dynamic login form
		c.itrisCSRFToken = el.ChildAttr("input[name=inloggen_csrf_protection]", "value")
		loginRequest["inloggen_csrf_protection"] = []byte(c.itrisCSRFToken)
		loginRequest["post0"] = []byte(el.ChildAttr("input[name=post0]", "value"))
	})

	// behave like a web browseer
	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"90\", \"Google Chrome\";v=\"90\"")
		r.Headers.Set("sec-ch-ua-mobile", "?0")
		r.Headers.Set("sec-fetch-site", "same-origin")
		r.Headers.Set("sec-fetch-mode", "navigate")
		r.Headers.Set("sec-fetch-user", "?1")
		r.Headers.Set("sec-fetch-dest", "document")
		r.Headers.Set("upgrade-insecure-requests", "1")
	})

	// parse login error
	var hasErrLogin error
	collector.OnResponse(func(resp *colly.Response) {
		if resp.StatusCode != http.StatusFound {
			hasErrLogin = checkLogin(string(resp.Body))
		}
	})

	// visit login page
	if err := collector.PostMultipart(loginURL, loginRequest); err != nil {
		return fmt.Errorf("itris connector: error while login: %w", err)
	}

	return hasErrLogin
}

func checkLogin(body string) error {
	errItrisLoginMsg := "Combinatie inlognaam / wachtwoord is niet bekend of onjuist. Controleer de invoer en probeer het opnieuw."
	errItrisBlockedMsg := "Om veiligheidsredenen is dit veld tijdelijk geblokkeerd, probeer het later nog eens"

	if strings.Contains(body, errItrisLoginMsg) {
		return connector.ErrAuthFailed
	} else if strings.Contains(body, errItrisBlockedMsg) {
		return ErrItrisBlocked
	}

	return ErrItrisLoginUnknown
}
