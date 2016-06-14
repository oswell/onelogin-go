package onelogin

import (
    "fmt"
    "errors"
)

const (
    ROLE_GET_ROLES         = "api/1/roles"
)

// Get all roles, or specify a name for a single role.
// Get_Roles("")           => Get all roles
// Get_Roles("operations") => Get operations role
func (o *OneLogin) Get_Roles(name string)([]OneLoginRole, error) {
    logger.Debugf("Requesting roles")

    oauth_token, err := o.Get_Token(); if err != nil {
        return nil, ErrorOcurred(err)
    }

    url := o.GetUrl(ROLE_GET_ROLES)

    headers := Headers(fmt.Sprintf("bearer:%s", oauth_token.Access_token))
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
        return -1, errors.New(fmt.Sprintf("Found %d roles matching %s, expected exactly 1.", len(roles), role))
    }

    return roles[0].Id, nil

}
