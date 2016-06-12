package onelogin

// OAuth token
type OneLogin_Token struct {
    Endpoint      string
    Client_id     string
    Client_secret string
    Access_token  string `json:"access_token"`
    Created_at    string `json:"created_at"`
    Expires_in    int    `json:"expires_in"`
    Refresh_token string `json:"refresh_token"`
}

// OAuth grant body
type OAuthGrantBody struct {
    Grant_type    string `json:"grant_type"`
    Refresh_token string `json:"refresh_token,omitempty"`
    Access_token  string `json:"access_token,omitempty"`
    Client_id     string `json:"client_id,omitempty"`
    Client_secret string `json:"client_secret,omitempty"`
}

// Object representing the JSON response for an OAuth authentication request.
type OAuthResponse struct {
    Status ResponseStatus `json:status`
    Data []struct {
        Access_token  string `json:"access_token"`
        Account_id    int    `json:"account_id"`
        Created_at    string `json:"created_at"`
        Expires_in    int    `json:"expires_in"`
        Refresh_token string `json:"refresh_token"`
        Token_type    string `json:"token_type"`
    } `json:"data"`
}
