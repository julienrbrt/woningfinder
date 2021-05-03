package itris

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

func (c *client) Login(username, password string) error {
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

	c.collector.OnHTML("form", func(el *colly.HTMLElement) {
		// fill dynamic login form
		loginRequest["inloggen_csrf_protection"] = []byte(el.ChildAttr("input[name=inloggen_csrf_protection]", "value"))
		loginRequest["post0"] = []byte(el.ChildAttr("input[name=post0]", "value"))
	})

	c.collector.OnResponse(func(resp *colly.Response) {
		fmt.Println(string(resp.Body))
	})

	// parse login page
	if err := c.collector.PostMultipart(loginURL, loginRequest); err != nil {
		return fmt.Errorf("error while login: %w", err)
	}

	return nil
}
