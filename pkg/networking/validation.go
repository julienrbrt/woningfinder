package networking

import (
	"fmt"
	"net/http"
)

// responseValidator checks the statusCode of the response and returns an error if the statusCode range is >= 200 & < 300 is not respected
func responseValidator(r *Response) error {
	if r.StatusCode >= http.StatusOK && r.StatusCode < http.StatusMultipleChoices {
		return nil
	}

	message := fmt.Sprintf(
		"Invalid response, expected the response code to be in the 200-299 range, got %d with body %s",
		r.StatusCode,
		getResponseBody(r),
	)

	return fmt.Errorf(message)
}

func getResponseBody(resp *Response) string {
	body, err := resp.CopyBody()
	if err != nil {
		return fmt.Sprintf("Unable to read the response body: %s", err)
	}

	return string(body)
}
