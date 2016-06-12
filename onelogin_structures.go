package onelogin

// The main OneLogin structure.
type OneLogin struct {
    // A hard coded base URL if you want to for some reason override the
    // defaults specified by shards
    CustomURL   string

    // OneLogin service shard (eu, us, etc)
    Shard         string
    Token         *OneLogin_Token

    OauthToken    string
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
