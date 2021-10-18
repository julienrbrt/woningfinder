package stripe

const (
	CheckoutSessionCompleted = "checkout.session.completed"
	InvoicePaid              = "invoice.paid"
	InvoicePaymentFailed     = "invoice.payment_failed"
)

func (c *client) WebhookSigningKey() string {
	return c.webookSigningKey
}
