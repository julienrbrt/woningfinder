package dewoonplaats

import "net/url"

// Host defines the De Woonplaats API domain
var Host = &url.URL{Scheme: "https", Host: "www.dewoonplaats.nl"}

const methodLogin = "Login"
const methodFind = "ZoekWoningen"

// request builds a De Woonplaats request
type request struct {
	ID     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}
