package email

import (
	"net/http"

	"github.com/mattevans/postmark-go"
)

type Client interface {
	Send(subject, htmlBody, to string) error
}

type client struct {
	postmark *postmark.Client
}

// NewClient permits to send an email
func NewClient(apiKey string) Client {
	return &client{
		postmark: postmark.NewClient(&http.Client{
			Transport: &postmark.AuthTransport{Token: apiKey},
		}),
	}
}
