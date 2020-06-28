package jexiasdkgo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client contains most data needed for each request
type Client struct {
	projectID   string
	projectZone string
	projectURL  string
	apiKey      string
	apiSecret   string
	token       Token
	http        *http.Client
}

// APITokenRequest is the JSON data sent to the /auth endpoint when authenticating with the API key
type APITokenRequest struct {
	Method string `json:"method"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

// UMSTokenRequest is the JSON data sent to the /auth endpoint when authenticating with user credentials
type UMSTokenRequest struct {
	Method string `json:"method"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Token is the response from the /auth request which contains the access and refresh tokens
// Currently, each token lasts 2 hours. https://docs.jexia.com/auth/#:~:text=token%20is%20valid%20for
type Token struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

// Option allows the client to be configured with different options.
type Option func(*Client)

// SetHTTPClient allows for a custom client to be set
func SetHTTPClient(http *http.Client) Option {
	return func(c *Client) {
		c.http = http
	}
}

// SetProjectURL allows for a custom url to be set which does not match that of the standard pattern
// Note: If such url contains strings such as the project ID, this needs to be computed before being passed through
func SetProjectURL(url string) Option {
	return func(c *Client) {
		c.projectURL = url
	}
}

// GetToken assigns the user token to the client for future use
// TODO: Trigger auto-refresh once called
func (c *Client) GetToken() {
	var token Token
	err := c.post(fmt.Sprintf("%v/auth", c.projectURL), APITokenRequest{
		Method: "apk",
		Key:    c.apiKey,
		Secret: c.apiSecret,
	}, &token)
	if err != nil {
		fmt.Printf("error from api. response: %v", err)
	}
	c.token = token
}

// NewClient is used to generate a new client for interacting with the API
func NewClient(id, zone, key, secret string, opts ...Option) *Client {
	client := &Client{
		projectID:   id,
		projectZone: zone,
		projectURL:  fmt.Sprintf("https://%v.%v.app.jexia.com", id, zone),
		apiKey:      key,
		apiSecret:   secret,
		token:       Token{},
		http:        &http.Client{},
	}
	for _, o := range opts {
		o(client)
	}
	return client
}

// buildRequest is used as a standard request where general headers are set
func (c *Client) buildRequest(method, url string, payload interface{}) (*http.Request, error) {

	b := bytes.NewBuffer(nil)
	if method == http.MethodPost || method == http.MethodPut {
		j, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		b = bytes.NewBuffer(j)
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}

	req.Header.Add("method", "apk")
	req.Header.Add("key", c.apiKey)
	req.Header.Add("secret", c.apiSecret)

	return req, err
}

// get returns an error if the http client cannot perform a HTTP GET for the provided URL.
func (c *Client) get(url string, target interface{}) error {
	req, err := c.buildRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("error from api")
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &target)
}

// save returns an error if the http client cannot correct the request
func (c *Client) save(httpMethod string, url string, payload interface{}, target interface{}) error {
	req, err := c.buildRequest(httpMethod, url, payload)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error from api. httpCode: %d, response: %s", resp.StatusCode, b)
	}
	return json.Unmarshal(b, &target)
}

// put returns an error if the http client cannot perform a HTTP PUT for the provided URL.
func (c *Client) put(url string, payload interface{}, target interface{}) error {
	return c.save(http.MethodPut, url, payload, target)
}

// post returns an error if the http client cannot perform a HTTP POST for the provided URL.
func (c *Client) post(url string, payload interface{}, target interface{}) error {
	return c.save(http.MethodPost, url, payload, target)
}
