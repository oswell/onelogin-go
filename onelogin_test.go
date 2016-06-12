package onelogin

import "testing"

// Test the GetUrl() method for generating request URLS
func TestGetUrl(t *testing.T) {

    shards := map[string]string{
        "us": "https://api.us.onelogin.com/",
        "eu": "https://api.eu.onelogin.com/",
    }

    for shard,url := range shards {
        o := OneLogin{Shard: shard}
        result := o.GetUrl("")
        if result != url {
            t.Errorf("GetUrl() for shard %s != %s (actual result was %s)", shard, url, result)
        }
    }
}

// Test the Get_Token() method for requesting the current or a new token.
func TestGet_Token(t *testing.T) {
    t.Errorf("No tests written yet for Get_Token()")
}

// Test the Authenticate() method for authenticating a user to OneLogin
func TestAuthenticate(t *testing.T) {
    t.Errorf("No tests written yet for Authenticate()")
}

// Test the VerifyToken() method for verifying an OTP (MFA) token.
func TestVerifyToken(t *testing.T) {
    t.Errorf("No tests written yet for VerifyToken()")
}
