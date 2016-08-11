package onelogin

import (
    "io"
    "bytes"
    "net/url"
    "io/ioutil"
    "net/http"
    "encoding/json"
)

/**
** client := HttpClient{Url:"http://www.google.com"}
** raw_response, err := client.go("GET", nil)
**
** raw_response, err := client.go("GET", &response_object)
**
**/

type HttpClient struct {
    Url     string                // Request URL (http://www.google.com/)
    Params  map[string]string     // Map of headers
    Headers map[string]string     // Map of query parameters
}

/**
 ** Perform an HTTP operation.
 **
 ** @reqObj Interface for passing in arbitrary objects that will be converted to JSON for transport
 ** @resObj Interface for returning arbitrary objects containing the unmarshaled response data
 **/
func (h *HttpClient) Request(method string, reqObj interface{}, respObj interface{})(interface{}, error) {
    body, err := h.EncodeBody(reqObj) ; if err != nil {
        return nil, ErrorOcurred(err)
    }

    response, err := h.Go(method, bytes.NewBufferString(string(body)), respObj) ; if err != nil {
        return nil, ErrorOcurred(err)
    }

    return response, nil
}
/**
 ** Perform the actual HTTP request, returning whatever response comes back as a byte array.
 **
 ** @method    The HTTP method to perform (GET, POST, PUT, DELETE, HEAD, OPTIONS, etc)
 ** @body      An io.Reader object containing the body to pass along with the request
 ** @respObj   The response object that we want data unmarshaled into, or nil if not needed.
 **/
func (h *HttpClient) Go(method string, body io.Reader, respObj interface{}) ([]byte, error) {
    base_url := h.CreateUrl(h.Url, h.Params)
    req, err := http.NewRequest(method, base_url, body) ; if err != nil {
        return nil, ErrorOcurred(err)
    }

    for k := range h.Headers {
        req.Header.Set(k, h.Headers[k])
    }

    client := &http.Client{}
    resp, err := client.Do(req) ; if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    response, err := ioutil.ReadAll(resp.Body) ; if err != nil {
        logger.Errorf("An error occurred while reading in the HTTP response body, %v", err)
        return nil, err
    }

    if respObj != nil {
        err = json.Unmarshal([]byte(response), &respObj) ; if err != nil {
            logger.Errorf("An error occurred unpacking the HTTP response body, %v", err)
            return nil, ErrorOcurred(err)
        }
    }

    logger.Debugf("Successfully called %s {%v}.", h.Url, h.Params)

    // body is now the response body, hopefully JSON.
    return response, nil
}

/**
 ** Encode the request body object to json.
 **
 ** requestBody   Interface for passing in arbitrary objects that will be converted to json for transport
 **/
func (h *HttpClient) EncodeBody(requestBody interface{})([]byte, error) {
    if requestBody == nil {
        return nil, nil
    }

    result, err := json.Marshal(requestBody) ; if err != nil {
        return nil, ErrorOcurred(err)
    }

    return result, nil
}

/**
 ** Generate a full request URL with provided parameters.
 **
 ** requrl      Base URL without required parameters.
 ** parameters  Map of key value pairs to add to the URL as query parameters.
 **/
func (h *HttpClient) CreateUrl(requrl string, parameters map[string]string)(string) {
    baseUrl, _ := url.Parse(requrl)
    params := url.Values{}
    for k, v := range parameters {
        params.Add(k, v)
    }
    baseUrl.RawQuery = params.Encode()
    return baseUrl.String()
}

/** Compile headers for the API call **/
func Headers(authorization string)(map[string]string) {
    return map[string]string{
        "Authorization": authorization,
        "Content-Type" : "application/json",
    }
}
