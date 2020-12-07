package dewoonplaats

import (
	"encoding/json"
	"fmt"
)

// response corresponds to a De Woonplaats response
type response struct {
	Err    interface{}     `json:"error"`
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result"`
}

func (r *response) Error() error {
	if r.Err != nil {
		return fmt.Errorf("de woonplaats error reponse: %v", r.Err.(string))
	}
	return nil
}
