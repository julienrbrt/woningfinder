package cryptocom

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// https://pay-docs.crypto.com/#api-reference-webhooks
const (
	PaymentCreated           = "payment.created"
	PaymentFulfill           = "payment.fulfill"
	PaymentCaptured          = "payment.captured"
	PaymentRefundRequest     = "payment.refund_requested"
	PaymentRefundTransferred = "payment.refund_transferred"

	timestampKey = "t"
	signatureKey = "v1"
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

// https://pay-docs.crypto.com/#api-reference-webhooks-webhook-signature
func (c *client) VerifyEvent(signatureHeader, eventRaw string) bool {
	parts := map[string]string{}
	for _, p := range strings.Split(signatureHeader, ",") {
		v := strings.Split(p, "=")
		if len(v) != 2 {
			continue
		}
		parts[v[0]] = v[1]
	}

	h := hmac.New(sha256.New, []byte(c.webookSigningKey))
	h.Write([]byte(fmt.Sprintf("%s.%s", parts[timestampKey], eventRaw)))

	return hex.EncodeToString(h.Sum(nil)) == parts[signatureKey]
}
