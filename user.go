package onelogin

// User Methods:
//
// func (o *OneLogin) Get_Users           (filter map[string]string)                                 (*[]OneLoginUser, error)
// func (o *OneLogin) Get_AppsForUser     (userId int)                                               ([]OneLoginApp,   error)
// func (o *OneLogin) Get_RolesForUser    (userId int)                                               ([]OneLoginRole,  error)
// func (o *OneLogin) Get_UsersWithRole   (role string)                                              ([]OneLoginUser,  error)
// func (o *OneLogin) Set_CustomAttribute (userid int, attributeName string, attributeValue string)  (error)


import (
    "fmt"
    "strconv"
)

const (
    USER_GET_USERS         = "api/1/users"
    USER_GET_APPS          = "api/1/users/%s/apps"
    USER_GET_USER_ROLES    = "api/1/users/%s/roles"

    USERS_SET_ATTRIBUTE    = "api/1/users/%s/set_custom_attributes"
    USER_CUSTOM_ATTRIBUTES = "api/1/users/custom_attributes"

    USER_AUTHENTICATE      = "api/1/login/auth"
    USER_VERIFY_FACTOR     = "api/1/login/verify_factor"
)

// Get_Users returns all users that match the filter.
// Allowed filters include:
//    directory_id, email, external_id, firstname, id, manager_ad_id,
//    role_id, samaccountname, since, until, username, userprincipalname
//
// filter can also include sort, fields, and limit.  See OneLogin documentation
// for more details (https://developers.onelogin.com/api-docs/1/users/get-users)
//
// OneLogin's API only returns at most 50 objects per call.  Until we enable paging
// on this end, we are going to page through and fetch all user objects.
//
// Returns a slice of OneLoginUser objects and an error if any error occurred.
func (o *OneLogin) Get_Users(filter map[string]string) (*[]OneLoginUser, error) {
    oauth_token, err := o.Get_Token(); if err != nil {
        return nil, ErrorOcurred(err)
    }

    url := o.GetUrl(USER_GET_USERS)
    headers := Headers(fmt.Sprintf("bearer:%s", oauth_token.Access_token))

    get_user_response := GetUserResponse{}

    client := HttpClient{Url: url, Headers: headers, Params: filter}
    _, err = client.Request("GET", nil, &get_user_response) ; if err != nil {
        logger.Errorf("An error occurred, %v", err)
        return nil, ErrorOcurred(err)
    }

    user_data := make([]OneLoginUser, 0)
    user_data = append(user_data, get_user_response.Data...)

    logger.Infof("Pagination link: %s", get_user_response.Pagination.After_cursor)

    next_cursor := get_user_response.Pagination.After_cursor
    for next_cursor != "" {

        filter["after_cursor"] = next_cursor
        client := HttpClient{Url: url, Headers: headers, Params: filter}

        get_user_response := GetUserResponse{}
        _, err = client.Request("GET", nil, &get_user_response) ; if err != nil {
            logger.Errorf("An error occurred, %v", err)
            return nil, ErrorOcurred(err)
        }
        next_cursor = get_user_response.Pagination.After_cursor

        logger.Infof("Pagination link: %s", next_cursor)
        user_data = append(user_data, get_user_response.Data...)
    }

    return &user_data, nil
}

// Get_AppsForUser returns all apps that a user has assigned.
// Returns a slice of OneLoginApp objects and an error if any errors occur
func (o *OneLogin) Get_AppsForUser(userId int) ([]OneLoginApp, error) {
    oauth_token, err := o.Get_Token(); if err != nil {
        fmt.Printf("ERROR ERROR ERROR")
        return nil, ErrorOcurred(err)
    }

    url := o.GetUrl(USER_GET_APPS, strconv.Itoa(userId))
    headers := Headers(fmt.Sprintf("bearer:%s", oauth_token.Access_token))

    params := map[string]string{}
    get_apps_for_user_response := GetUserAppsResponse{}

    client := HttpClient{Url: url, Headers: headers, Params: params}
    _, err = client.Request("GET", nil, &get_apps_for_user_response) ; if err != nil {
        logger.Errorf("Error requesting apps for user %d, %v", userId, err)
        return nil, ErrorOcurred(err)
    }

    return get_apps_for_user_response.Data, nil
}

// Get_RolesForUser returns all apps that a user has assigned.
// Returns a slice of integers representing role ids and an error if any errors occur
func (o *OneLogin) Get_RolesForUser(userId int) ([]int, error) {
    oauth_token, err := o.Get_Token(); if err != nil {
        return nil, ErrorOcurred(err)
    }

    url := o.GetUrl(USER_GET_USER_ROLES, strconv.Itoa(userId))
    headers := Headers(fmt.Sprintf("bearer:%s", oauth_token.Access_token))

    params := map[string]string{}
    get_roles_for_user_response := GetUserRolesResponse{}

    client := HttpClient{Url: url, Headers: headers, Params: params}
    _, err = client.Request("GET", nil, &get_roles_for_user_response) ; if err != nil {
        return nil, ErrorOcurred(err)
    }

    return get_roles_for_user_response.Data, nil
}

// Get_UsersWithRole returns all users that have the specified role.  If no role
// is given (empty string) then all users will be returned.
// Returns a slice of OneLoginUser objects and an error if any error occurred.
func (o *OneLogin) Get_UsersWithRole(role string)([]OneLoginUser, error) {

    oauth_token, err := o.Get_Token(); if err != nil {
        return nil, ErrorOcurred(err)
    }

    url := o.GetUrl(USER_GET_USERS)
    headers := Headers(fmt.Sprintf("bearer:%s", oauth_token.Access_token))

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

// Set_CustomAttribute sets the custom attribute (attributeName) with the value (attributeValue)
// for the user with the given ID.
// Returns an error if any error occurred.
func (o *OneLogin) Set_CustomAttribute(userid int, attributeName string, attributeValue string)(error) {
    logger.Debugf("Setting custom attribute %s for user id %d", attributeName, userid)

    oauth_token, err := o.Get_Token(); if err != nil {
        return ErrorOcurred(err)
    }

    url := o.GetUrl(USERS_SET_ATTRIBUTE, strconv.Itoa(userid))
    headers := Headers(fmt.Sprintf("bearer:%s", oauth_token.Access_token))
    client := HttpClient{Url: url, Headers: headers}

    attributes := SetCustomAttributeRequest{Custom_Attributes: map[string]string{ attributeName: attributeValue }}
    _, err = client.Request("PUT", &attributes, nil) ; if err != nil {
        return ErrorOcurred(err)
    }

    return nil
}
