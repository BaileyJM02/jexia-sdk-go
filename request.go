package jexiasdkgo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// requestOption allows the request builder to be configured with different options.
type requestOption func(*http.Request)

// addHeader allows individual headers to be set if and when required
func addHeader(key, value string) requestOption {
	return func(r *http.Request) {
		r.Header.Add(key, value)
	}
}

// setBody sets the body of the request when we are making calls such as a http post
func setBody(body []byte) requestOption {
	return func(r *http.Request) {
		r.Body = ioutil.NopCloser(bytes.NewReader(body))
		r.ContentLength = int64(len(body))
		// We do not set the TransferEncoding encoding at the moment as we are not aware of an edge case requiring it
		// r.TransferEncoding = ""
	}
}

// useAPKMethod is a short-hand function for making calls using the APK Jexia method
func useAPKMethod(key, secret string) requestOption {
	return func(r *http.Request) {
		r.Header.Add("Method", "apk")
		r.Header.Add("Key", key)
		r.Header.Add("Secret", secret)
	}
}

// addToken is a short-hand function for setting the 'Authorization' header
func addToken(accessToken string) requestOption {
	return func(r *http.Request) {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))
	}
}

// buildRequest is used to build the request by sending requestOption functions
func (c *Client) buildRequest(method, url string, opts ...requestOption) (*http.Request, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(nil))
	if err != nil {
		return nil, err
	}

	if method == http.MethodPost || method == http.MethodPut {
		// Default headers for these http types, will be overridden in an option function
		req.Header.Add("Content-Type", "application/json")
	}

	for _, o := range opts {
		o(req)
	}

	return req, err
}

// executeRequest calls the http.Do function
func (c *Client) executeRequest(req *http.Request, target interface{}) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return &Error{
			ID:        "e005",
			Message:   fmt.Errorf("Unable to execute http request: %w", err).Error(),
			Origin:    Internal,
			Temporary: false,
		}
	}

	defer resp.Body.Close()

	err = checkForAPIError(resp)
	if err != nil {
		return err
	}

	b, err := read(resp.Body)
	if err != nil {
		return err
	}
	return unmarshal(b, &target)
}

// get performs a http get request
func (c *Client) get(url string, target interface{}, opts ...requestOption) error {
	req, err := c.buildRequest(http.MethodGet, url, opts...)
	if err != nil {
		return err
	}
	return c.executeRequest(req, target)
}

// put performs a http put request for updating data
func (c *Client) put(url string, target interface{}, opts ...requestOption) error {
	req, err := c.buildRequest(http.MethodPut, url, opts...)
	if err != nil {
		return err
	}
	return c.executeRequest(req, target)
}

// post performs a http post request for passing data to the endpoint
func (c *Client) post(url string, target interface{}, opts ...requestOption) error {
	req, err := c.buildRequest(http.MethodPost, url, opts...)
	if err != nil {
		return err
	}
	return c.executeRequest(req, target)
}
