package stripe

const (
	CheckoutSessionCompleted    = "checkout.session.completed"
	InvoicePaid                 = "invoice.paid"
	InvoicePaymentFailed        = "invoice.payment_failed"
	CustomerSubscriptionDeleted = "customer.subscription.deleted"
)

func (c *client) WebhookSigningKey() string {
	return c.webookSigningKey
}
