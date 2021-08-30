package stripe

func (c *client) WebhookSigningKey() string {
	return c.webookSigningKey
}
