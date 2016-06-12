package onelogin

import (
    "fmt"
    "errors"
)

const (
    OAUTH_ENDPOINT     = "auth/oauth2"
)

// Fetch an OAuth token from the OneLogin API to be used for resource authentication.
func (t *OneLogin_Token) Get() (error) {
    headers := Headers(fmt.Sprintf("client_id:%s,client_secret:%s", t.Client_id, t.Client_secret))
    body := OAuthGrantBody{Grant_type: "client_credentials"}
    response := OAuthResponse{}

    url := fmt.Sprintf("%s/%s/token", t.Endpoint, OAUTH_ENDPOINT)
    client := HttpClient{Url: url, Headers: headers}

    _, err := client.Request("POST", &body, &response) ; if err != nil {
        logger.Errorf("Error requesting token, %v", err)
        return ErrorOcurred(err)
    }

    if response.Status.Error {
        err = errors.New(response.Status.Message)
        logger.Errorf("Error requesting token, %v", err)
        return ErrorOcurred(err)
    }

    // Store the relevant data for later consumption.
    t.Access_token = response.Data[0].Access_token
    t.Created_at = response.Data[0].Created_at
    t.Expires_in = response.Data[0].Expires_in
    t.Refresh_token = response.Data[0].Refresh_token

    return nil
}

// Refresh an OAuth token from the OneLogin API to be used for resource authentication.
func (t *OneLogin_Token) Refresh() (error) {
    body := OAuthGrantBody{
        Grant_type: "client_credentials",
        Refresh_token: t.Refresh_token,
        Access_token: t.Access_token,
    }
    response := OAuthResponse{}

    url := fmt.Sprintf("%s/%s/token", t.Endpoint, OAUTH_ENDPOINT)
    client := HttpClient{Url: url}

    _, err := client.Request("POST", &body, &response) ; if err != nil {
        return ErrorOcurred(err)
    }

    if response.Status.Error {
        return ErrorOcurred(errors.New(response.Status.Message))
    }

    // Store the relevant data for later consumption.
    t.Access_token = response.Data[0].Access_token
    t.Created_at = response.Data[0].Created_at
    t.Expires_in = response.Data[0].Expires_in
    t.Refresh_token = response.Data[0].Refresh_token

    return nil
}

// Revoke an OAuth token from the OneLogin API.
func (t *OneLogin_Token) Revoke() (error) {
    body := OAuthGrantBody{
        Grant_type: "client_credentials",
        Access_token: t.Access_token,
        Client_id: t.Client_id,
        Client_secret: t.Client_secret,
    }
    response := OAuthResponse{}

    url := fmt.Sprintf("%s/%s/token", t.Endpoint, OAUTH_ENDPOINT)
    client := HttpClient{Url: url}

    _, err := client.Request("POST", &body, &response) ; if err != nil {
        return ErrorOcurred(err)
    }

    if response.Status.Error {
        return ErrorOcurred(errors.New(response.Status.Message))
    }

    // Store the relevant data for later consumption.
    t.Access_token = response.Data[0].Access_token
    t.Created_at = response.Data[0].Created_at
    t.Expires_in = response.Data[0].Expires_in
    t.Refresh_token = response.Data[0].Refresh_token

    return nil
}
