package onelogin

import (
    "fmt"
    "errors"
    "strconv"
)

type OneLogin struct {
    OauthToken    string
    SubDomain     string
    Client_id     string
    Client_secret string
}

const ONELOGIN_URL = "https://api.us.onelogin.com"

const (
    OAUTH_AUTHENTICATE     = "auth/oauth2/token"
    USER_AUTHENTICATE      = "api/1/login/auth"
    USER_CUSTOM_ATTRIBUTES = "api/1/users/custom_attributes"
    USER_VERIFY_FACTOR     = "api/1/login/verify_factor"
    USER_GET_USERS         = "api/1/users"
    USERS_SET_ATTRIBUTE    = "api/1/users/%d/set_custom_attributes"
    ROLE_GET_ROLES         = "api/1/roles"
)

func (o *OneLogin) createUrl(uri string)(string) {
    return fmt.Sprintf("%s/%s", ONELOGIN_URL, uri)
}

/**
 ** Get the user object
 **/
func (o *OneLogin) Get_UserByUsername(username string)(*OneLoginUser, error) {
    logger.Debugf("Requesting user data for %s", username)

    oauth_token, err := o.Get_OAuth_Token(); if err != nil {
        return nil, ErrorOcurred(err)
    }

    url := o.createUrl(USER_GET_USERS)
    headers := o.headers(fmt.Sprintf("bearer:%s", oauth_token))
    params := map[string]string{ "username": username }

    get_user_response := GetUserResponse{}

    client := HttpClient{Url: url, Headers: headers, Params: params}
    _, err = client.Request("GET", nil, &get_user_response) ; if err != nil {
        return nil, ErrorOcurred(err)
    }

    if len(get_user_response.Data) != 1 {
        logger.Debugf("Found %d users, but expected exactly 1.", len(get_user_response.Data))
        return nil, ErrorOcurred(errors.New("No user found."))
    }

    return &get_user_response.Data[0], nil
}

// Get all roles, or specify a name for a single role.
// Get_Roles("")           => Get all roles
// Get_Roles("operations") => Get operations role
func (o *OneLogin) Get_Roles(name string)([]OneLoginRole, error) {
    logger.Debugf("Requesting roles")

    oauth_token, err := o.Get_OAuth_Token(); if err != nil {
        return nil, ErrorOcurred(err)
    }

    url := o.createUrl(ROLE_GET_ROLES)
    headers := o.headers(fmt.Sprintf("bearer:%s", oauth_token))
    params := map[string]string{}

    if name != "" {
        params = map[string]string{ "name": name }
    }

    get_role_response := GetRoleResponse{}

    client := HttpClient{Url: url, Headers: headers, Params: params}
    _, err = client.Request("GET", nil, &get_role_response) ; if err != nil {
        return nil, ErrorOcurred(err)
    }

    roles := get_role_response.Data

    return roles, nil
}

func (o *OneLogin) Get_Role_Id(role string)(int, error) {
    logger.Debugf("Requesting role %s", role)

    roles, err := o.Get_Roles(role) ; if err != nil {
        return -1, ErrorOcurred(err)
    }
    if len(roles) != 1 {
        return -1, errors.New(fmt.Sprintf("Found %d roles matching %s, expected only 1.", len(roles), role))
    }

    return roles[0].Id, nil

}

func (o *OneLogin) Get_UsersWithRole(role string)([]OneLoginUser, error) {

    oauth_token, err := o.Get_OAuth_Token(); if err != nil {
        return nil, ErrorOcurred(err)
    }

    url := o.createUrl(fmt.Sprintf(USER_GET_USERS))
    headers := o.headers(fmt.Sprintf("bearer:%s", oauth_token))

    params := map[string]string{}

    // if role is specified, add the query parameter for it.
    if role != "" {
        logger.Debugf("Requesting all users in role %s", role)

        role_id, err := o.Get_Role_Id(role) ; if err != nil {
            return nil, ErrorOcurred(err)
        }

        params = map[string]string{ "role_id": fmt.Sprintf("%d", role_id) }
    }

    get_user_response := GetUserResponse{}

    client := HttpClient{Url: url, Headers: headers, Params: params}
    _, err = client.Request("GET", nil, &get_user_response) ; if err != nil {
        return nil, ErrorOcurred(err)
    }

    return get_user_response.Data, nil
}

func (o *OneLogin) Set_CustomAttribute(userid int, attributeName string, attributeValue string)(error) {
    logger.Debugf("Setting custom attribute %s for user id %d", attributeName, userid)

    oauth_token, err := o.Get_OAuth_Token(); if err != nil {
        return ErrorOcurred(err)
    }

    url := o.createUrl(fmt.Sprintf(USERS_SET_ATTRIBUTE, userid))
    headers := o.headers(fmt.Sprintf("bearer:%s", oauth_token))
    client := HttpClient{Url: url, Headers: headers}

    attributes := SetCustomAttributeRequest{Custom_Attributes: map[string]string{ attributeName: attributeValue }}
    _, err = client.Request("PUT", &attributes, nil) ; if err != nil {
        return ErrorOcurred(err)
    }

    return nil
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

    oauth_token, err := o.Get_OAuth_Token(); if err != nil {
        return ErrorOcurred(err)
    }

    url := o.createUrl(USER_AUTHENTICATE)
    headers := o.headers(fmt.Sprintf("bearer:%s", oauth_token))

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

    oauth_token, err := o.Get_OAuth_Token(); if err != nil {
        return ErrorOcurred(err)
    }

    url := o.createUrl(USER_VERIFY_FACTOR)
    headers := o.headers(fmt.Sprintf("bearer:%s", oauth_token))
    client := HttpClient{Url: url, Headers: headers}

    _, err = client.Request("POST", &verify_request, &auth_response); if err != nil {
        return ErrorOcurred(err)
    }

    if auth_response.Status.Error {
        return ErrorOcurred(errors.New(auth_response.Status.Message))
    }

    return nil
}

/**
 ** Perform OAuth authentication to retrieve a token for further API calls.
 **/
func (o *OneLogin) AuthenticateOauth(client_id string, client_secret string)(error) {
    logger.Debug("Retrieving oauth token.")

    headers := o.headers(fmt.Sprintf("client_id:%s,client_secret:%s", client_id, client_secret))
    body := OAuthGrantBody{Grant_Type: "client_credentials"}
    response := OAuthResponse{}

    client := HttpClient{Url: o.createUrl(OAUTH_AUTHENTICATE), Headers: headers}
    _, err := client.Request("POST", &body, &response) ; if err != nil {
        return ErrorOcurred(err)
    }

    if response.Status.Error {
        return ErrorOcurred(errors.New(response.Status.Message))
    }

    o.OauthToken = response.Data[0].Access_token
    return nil
}

func (o *OneLogin) Get_OAuth_Token()(string, error) {
    if o.OauthToken == "" {

        // Fetch configuration values for the OAuth client_id and secret.
        client_id := o.Client_id
        client_secret := o.Client_secret

        err := o.AuthenticateOauth(client_id, client_secret) ; if err != nil {
            return "", err
        }

        logger.Debugf("Successfully fetched an OAuth Token.")
    }

    return o.OauthToken, nil
}

/** Compile headers for the API call **/
func (o *OneLogin) headers(authorization string)(map[string]string) {
    return map[string]string{
        "Authorization": authorization,
        "Content-Type" : "application/json",
    }
}
