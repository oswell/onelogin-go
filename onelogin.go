package onelogin

import (
    "fmt"
    "errors"
    "strconv"
    "github.com/op/go-logging"
)

func New(shard string, client_id string, client_secret string, subdomain string, loglevel logging.Level)(*OneLogin) {
    ol := OneLogin{Shard:shard, Client_id: client_id, Client_secret:client_secret, SubDomain: subdomain}
    ol.SetLogLevel(loglevel)
    return &ol
}

func (o *OneLogin) SetLogLevel(loglevel logging.Level) {
    SetLogLevel(loglevel)
}

// GetUrl creates a URL given the URI and any given args.
// Returns a URL
func (o *OneLogin) GetUrl(uri string, args... string)(string) {
    // Handle cases where the uri requires variable replacements (ie.  /api/1/user/%d/roles)
    fulluri := uri

    if len(args) > 0 {
        // Convert to slice of interface so that the slice can be sent in as a variadic argument
        argint := make([]interface{}, len(args))
        for index, value := range args { argint[index] = value }
        fulluri = fmt.Sprintf(uri, argint...)
    }

    if o.CustomURL != "" {
        return fmt.Sprintf("%s/%s", o.CustomURL, fulluri)
    }

    return fmt.Sprintf("%s/%s", fmt.Sprintf(ONELOGIN_URL, o.Shard), fulluri)
}

// Convenience function to always return a token, generating or refreshing if necessary.
// TODO: Enable refreshing when necessary.
//
func (o *OneLogin) Get_Token()(*OneLogin_Token, error) {
    if o.Token == nil {
        o.Token = &OneLogin_Token{
            Endpoint     : o.GetUrl(""),
            Client_id    : o.Client_id,
            Client_secret: o.Client_secret,
        }
        err := o.Token.Get() ; if err != nil {
            return nil, ErrorOcurred(err)
        }
    }

    logger.Debugf("Token: %s", o.Token)
    return o.Token, nil
}


/**
 ** Performs an authentication request for a user.  An optional MFA token can be passed in with
 ** the assumption that it will be required.  If no token is passed in and one is required,
 ** authentication will fail.
 **/
func (o *OneLogin) Authenticate(username string, password string, token string)(error) {
    logger.Debugf("Authenticating user %s", username)

    auth_request := AuthenticationRequest{
        Username_or_email: username,
        Password: password,
        Subdomain: o.SubDomain,
    }

    auth_response := AuthResponse{}

    oauth_token, err := o.Get_Token(); if err != nil {
        return ErrorOcurred(err)
    }

    url := o.GetUrl(USER_AUTHENTICATE)
    headers := Headers(fmt.Sprintf("bearer:%s", oauth_token))

    client := HttpClient{Url: url, Headers: headers}
    _, err = client.Request("POST", &auth_request, &auth_response) ; if err != nil {
        return ErrorOcurred(err)
    }
    if auth_response.Status.Error {
        logger.Errorf("Error authenticating user: %s", auth_response.Status.Message)
        return ErrorOcurred(errors.New("An error occurred while authenticating."))
    }

    /** TODO: do not assume MFA is required. **/
    return o.VerifyToken(
        strconv.Itoa(auth_response.Data[0].Devices[0].Device_id),
        auth_response.Data[0].State_token,
        token,
    )
}

/**
 ** Verify an MFA token
 **/
func (o *OneLogin) VerifyToken(device_id string, state_token string, token string)(error) {
    logger.Debugf("Verifying MFA token for device %s", device_id)
    verify_request := VerifyTokenRequest{
        Device_id   : device_id,
        State_token : state_token,
        Otp_token   : token,
    }

    auth_response := &AuthResponse{}

    oauth_token, err := o.Get_Token(); if err != nil {
        return ErrorOcurred(err)
    }

    url := o.GetUrl(USER_VERIFY_FACTOR)
    headers := Headers(fmt.Sprintf("bearer:%s", oauth_token))
    client := HttpClient{Url: url, Headers: headers}

    _, err = client.Request("POST", &verify_request, &auth_response); if err != nil {
        return ErrorOcurred(err)
    }

    if auth_response.Status.Error {
        return ErrorOcurred(errors.New(auth_response.Status.Message))
    }

    return nil
}
