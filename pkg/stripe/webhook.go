package stripe

const (
	PaymentIntentSucceeded = "payment_intent.succeeded"
)

func (c *client) WebhookSigningKey() string {
	return c.webookSigningKey
}
