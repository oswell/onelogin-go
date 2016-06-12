package onelogin


import (
    "fmt"
    "testing"
    "strings"
    "net/http"
    "net/http/httptest"
    "io/ioutil"
)

var TestDataMap = map[string]string{
    "/auth/oauth2/token"            : "testdata/get_token.js",
    USER_GET_USERS                  : "testdata/get_users.js",
    fmt.Sprintf(USER_GET_APPS, "1") : "testdata/get_apps_for_user.js",
    ROLE_GET_ROLES                  : "testdata/get_roles.js",
}

var Main_Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    if strings.HasPrefix(path, "/") {
        path = path[1:]
    }
    filename := TestDataMap[path]

    logger.Infof("Test data filename for %s: %s\n", path, filename)
    logger.Infof("%s", fmt.Sprintf(USER_GET_APPS, "1"))

    // If no response data is found, return an empty json blob.
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        w.WriteHeader(502)
        fmt.Fprintf(w, string("{}"))
        return
    }

    w.WriteHeader(200)
	fmt.Fprintf(w, string(data))
})

func TestGet_Users(t *testing.T) {
    ts := httptest.NewServer(Main_Handler)
    o := new(OneLogin)
    o.CustomURL = ts.URL
    defer ts.Close()

    result, err := o.Get_Users(nil)

    if err != nil {
        t.Errorf("Failed to get users, %v", err)
        return
    }

    // Ensure length.
    if len(*result) != 2 {
        t.Errorf("Expected 2 users returned, found %d", len(*result))
        return
    }

    if (*result)[0].Username != "hzhang" {
        t.Errorf("Expected the first username to be hzhang")
        return
    }
}

func TestGet_AppsForUser(t *testing.T) {
    ts := httptest.NewServer(Main_Handler)
    o := new(OneLogin)
    o.CustomURL = ts.URL
    defer ts.Close()

    result, err := o.Get_AppsForUser(1)

    if err != nil {
        t.Errorf("Failed to get apps, %v", err)
        return
    }

    // Ensure length.
    if len(result) != 3 {
        t.Errorf("Expected 3 apps returned, found %d", len(result))
        return
    }

    if (result)[0].Login_id != 66666666 {
        t.Errorf("Expected the login ID to be 66666666")
        return
    }
}

func TestGet_RolesForUser(t *testing.T) {
    ts := httptest.NewServer(Main_Handler)
    o := new(OneLogin)
    o.CustomURL = ts.URL
    defer ts.Close()

    _, err := o.Get_RolesForUser(1)

    if err != nil {
        t.Errorf("Failed to get roles, %v", err)
        return
    }
}

func TestGet_UsersWithRole(t *testing.T) {
    ts := httptest.NewServer(Main_Handler)
    o := new(OneLogin)
    o.CustomURL = ts.URL
    defer ts.Close()

    _, err := o.Get_UsersWithRole("marketing")

    if err != nil {
        t.Errorf("Failed to get roles, %v", err)
        return
    }
}

func TestSet_CustomAttribute(t *testing.T) {
    t.Errorf("No tests written yet for Set_CustomAttribute()")
}
