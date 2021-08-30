package cryptocom

// https://pay-docs.crypto.com/#api-reference-webhooks
const (
	PaymentCreated           = "payment.created"
	PaymentFulfill           = "payment.fulfill"
	PaymentCaptured          = "payment.captured"
	PaymentRefundRequest     = "payment.refund_requested"
	PaymentRefundTransferred = "payment.refund_transferred"
)

type WebhookEvent struct {
	ID         string `json:"id"`
	ObjectType string `json:"object_type"`
	Type       string `json:"type"`
	Created    int    `json:"created"`
	Data       struct {
		Object CryptoCheckoutSession `json:"object"`
	} `json:"data"`
}
