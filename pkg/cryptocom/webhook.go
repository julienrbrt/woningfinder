package cryptocom

// https://pay-docs.crypto.com/#api-reference-webhooks
const (
	PaymentCreated           = "payment.created"
	PaymentFulfill           = "payment.fulfill"
	PaymentCaptured          = "payment.captured"
	PaymentRefundRequest     = "payment.refund_requested"
	PaymentRefundTransferred = "payment.refund_transferred"
)
