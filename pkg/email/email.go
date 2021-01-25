package email

type Client interface {
	Send(subject, html, plain, to string) error
}

type client struct {
	name     string
	from     string
	password string
	server   string
	port     int
}

// NewClient permits to send an email
func NewClient(name, from, password, server string, port int) Client {
	return &client{name, from, password, server, port}
}
