package networking

import (
	"fmt"
	"net/http"
)

// responseValidator checks the statusCode of the response and returns an error if the status code not in 200-299 range
func responseValidator(r *Response) error {
	if r.StatusCode >= http.StatusOK && r.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	message := fmt.Sprintf(
		"Invalid response, expected the response code to be in the 200-299 range, got %d with body %s",
		r.StatusCode,
		getResponseBody(r),
	)

	err := NewNetworkingError(r, message)
	if isTemporary(r.StatusCode) {
		return NewTemporaryNetError(err)
	}

	return err
}

func isTemporary(status int) bool {
	codes := []int{
		http.StatusRequestTimeout,
		http.StatusTooManyRequests,
		http.StatusLocked,
		http.StatusTooEarly,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout,
		http.StatusInsufficientStorage,
	}
	for _, k := range codes {
		if k == status {
			return true
		}
	}

	return false
}

func getResponseBody(resp *Response) string {
	body, err := resp.CopyBody()
	if err != nil {
		return fmt.Sprintf("Unable to read the response body: %s", err)
	}

	return string(body)
}
