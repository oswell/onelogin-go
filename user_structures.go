package onelogin

// Object representing a single App description
type OneLoginApp struct {
    Id          int    `json:"id"`
    Name        string `json:"name"`
    Icon        string `json:"icon"`
    Provisioned string `json:"provisioned"`
    Extension   bool   `json:"extension"`
    Login_id    int32  `json:"login_id"`
    Personal    bool   `json:"personal"`
}

// Object representing an HTTP request body for verifying an OTP token.
type VerifyTokenRequest struct {
    Device_id    string `json:"device_id"`
    State_token  string `json:"state_token"`
    Otp_token    string `json:"otp_token"`
}

// Object representing an HTTP request body for verifying user credentials.
type AuthenticationRequest struct {
    Username_or_email string `json:"username_or_email"`
    Password          string `json:"password"`
    Subdomain         string `json:"subdomain"`
}

type SetCustomAttributeRequest struct {
    Custom_Attributes map[string]string `json:"custom_attributes"`
}

// Object representing the pagination data for a request.
type ResponsePagination struct {
    Before_cursor  string `json:"before_cursor"`
    After_cursor   string `json:"after_cursor"`
    Previous_link  string `json:"previous_link"`
    Next_link      string `json:"next_link"`
}

// Object representing the JSON response for a user authentication
// request.
type AuthResponse struct {
    Status ResponseStatus `json:"status"`
    Data []struct {
        State_token    string `json:"state_token"`
        Devices   []struct {
            Device_type string `json:"device_type"`
            Device_id   int    `json:"device_id"`
        }  `json:devices`
        User        struct {
            Lastname  string `json:"lastname"`
            Username  string `json:"username"`
            Email     string `json:"email"`
            Id        int    `json:"id"`
            Firstname string `json:"firstname"`
        } `json:user`
        Callback_url  string `json:"callback_url"`
    } `json:data`
}

// Object representing user data.
type OneLoginUser struct {
    Activated_at           string            `json:"activated_at"`
    Created_at             string            `json:"created_at"`
    Email                  string            `json:"email"`
    Username               string            `json:"username"`
    Firstname              string            `json:"firstname"`
    Group_id               int               `json:"group_id"`
    Id                     int               `json:"id"`
    Invalid_login_attempts int               `json:"invalid_login_attempts"`
    Invitation_sent_at     string            `json:"invitation_sent_at"`
    Last_login             string            `json:"last_login"`
    Lastname               string            `json:"lastname"`
    Locked_until           string            `json:"locked_until"`
    Notes                  string            `json:"notes"`
    Openid_name            string            `json:"openid_name"`
    Locale_code            string            `json:"locale_code"`
    Password_changed_at    string            `json:"password_changed_at"`
    Phone                  string            `json:"phone"`
    Status                 int               `json:"status"`
    Updated_at             string            `json:"updated_at"`
    Distinguished_name     string            `json:"distinguished_name"`
    External_id            int               `json:'external_id'`
    Directory_id           int               `json:"directory_id"`
    Member_of              []string          `json:"member_of"`
    Samaccountname         string            `json:"samaccountname"`
    Userprincipalname      string            `json:"userprincipalname"`
    Manager_ad_id          int               `json:"manager_ad_id"`
    Role_id                []int             `json:"role_id"`
    Custom_attributes      map[string]string `json:"custom_attributes"`
}

type OneLoginRole struct {
    Id   int    `json:"id"`
    Name string `json:"name,omitempty"`
}

// Object representing the result of a get user request.
type GetUserResponse struct {
    Status     ResponseStatus     `json:"status"`
    Pagination ResponsePagination `json:"pagination"`
    Data       []OneLoginUser     `json:"data"`
}

// Object representing the result of a get user apps request.
type GetUserAppsResponse struct {
    Status     ResponseStatus     `json:"status"`
    Pagination ResponsePagination `json:"pagination"`
    Data       []OneLoginApp      `json:"data"`
}
// Object representing the result of a get user roles request.
type GetUserRolesResponse struct {
    Status     ResponseStatus     `json:"status"`
    Pagination ResponsePagination `json:"pagination"`
    Data       []int              `json:"data"`
}

// Object representing the result of a get user request.
type GetRoleResponse struct {
    Status     ResponseStatus     `json:"status"`
    Pagination ResponsePagination `json:"pagination"`
    Data       []OneLoginRole     `json:"data"`
}
