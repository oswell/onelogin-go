package onelogin

// The main OneLogin structure.
type OneLogin struct {
    // Allow custom base URL to override the generated URL
	CustomURL string

	// OneLogin service shard (eu, us, etc)
	Shard string

	// Token struct for managing the OAuth token
	Token *OneLogin_Token

    SubDomain     string
    Client_id     string
    Client_secret string
}

// Object representing the JSON "status" portion of a request response.
// All responses contain a "status" object.
type ResponseStatus struct {
    Code    int    `json:"code"`
    Error   bool   `json:"error"`
    Message string `json:"message"`
    Type    string `json:"type"`
}

const ONELOGIN_URL = "https://api.%s.onelogin.com"
