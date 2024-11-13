package email

type Client interface {
	Send(subject, htmlBody, to string) error
}

type client struct {
	from     string
	username string
	password string
	server   string
	port     int
}

// NewClient permits to send an email
func NewClient(from, username, password, server string, port int) Client {
	return &client{from, username, password, server, port}
}
