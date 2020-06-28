package jexiasdkgo

import (
	"fmt"
	"net/http"
	"time"
)

// Client contains most data needed for each request
type Client struct {
	projectID    string
	projectZone  string
	projectURL   string
	token        Token
	tokenRequest interface{}
	http         *http.Client
}

// APKTokenRequest is the JSON data sent to the /auth endpoint when authenticating with the API key
type APKTokenRequest struct {
	Method string `json:"method"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

// UMSTokenRequest is the JSON data sent to the /auth endpoint when authenticating with user credentials
type UMSTokenRequest struct {
	Method   string `json:"method"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Token is the response from the /auth request which contains the access and refresh tokens
// Currently, each token lasts 2 hours. https://docs.jexia.com/auth/#:~:text=token%20is%20valid%20for
type Token struct {
	Access   string `json:"access_token"`
	Refresh  string `json:"refresh_token"`
	Lifetime time.Duration
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

// SetToken assigns the user token to the client for future use
func (c *Client) SetToken(token Token) {
	c.token = token
}

// Token assigns the user token to the client for future use
func (c *Client) fetchToken(target *Token) error {
	payload, _ := marshal(c.tokenRequest)
	err := c.post(
		fmt.Sprintf("%v/auth", c.projectURL),
		&target,
		setBody(payload),
	)
	if err != nil {
		fmt.Printf("error from api. response: %v", err)
		return err
	}
	return nil
}

// SetTokenLifetime sets the duration before a token refresh is called
// Note: This currently only applies after the first 118 minute loop
// TODO: Ensure that this new duration is set immediately and not after the current loop
func (c *Client) SetTokenLifetime(duration time.Duration) {
	c.token.Lifetime = duration
}

// UseAPKToken assigns the user token to the client for future use
func (c *Client) UseAPKToken(apiKey, apiSecret string) {
	var token Token
	c.tokenRequest = APKTokenRequest{
		Method: "apk",
		Key:    apiKey,
		Secret: apiSecret,
	}
	err := c.fetchToken(&token)
	if err != nil {
		fmt.Printf("Unable to fetch token")
	}
	c.SetToken(token)
	// 2 hours minus 2 minutes to ensure we never lose the token
	c.SetTokenLifetime(118 * time.Minute)
	go c.startRefreshing()
}

// UseUMSToken assigns the user token to the client for future use
func (c *Client) UseUMSToken(email, password string) {
	var token Token
	c.tokenRequest = UMSTokenRequest{
		Method:   "ums",
		Email:    email,
		Password: password,
	}
	err := c.fetchToken(&token)
	if err != nil {
		fmt.Printf("Unable to fetch token")
	}
	c.SetToken(token)
	// 2 hours minus 2 minutes to ensure we never lose the token
	c.SetTokenLifetime(118 * time.Minute)
	go c.startRefreshing()
}

// RefreshToken triggers a token refresh once called
func (c *Client) RefreshToken() {
	var token Token
	payload, _ := marshal(c.tokenRequest)
	err := c.post(fmt.Sprintf("%v/auth/refresh", c.projectURL), &token, setBody(payload), addToken(c.token.Access))
	if err != nil {
		fmt.Printf("error from api. response: %v", err)
	}

	// Pass the new refresh token over to the client
	c.token.Refresh = token.Refresh
}

func (c *Client) startRefreshing() {
	for {
		time.Sleep(c.token.Lifetime)
		go c.RefreshToken()
	}
}

// NewClient is used to generate a new client for interacting with the API
func NewClient(id, zone string, opts ...Option) *Client {
	client := &Client{
		projectID:    id,
		projectZone:  zone,
		projectURL:   fmt.Sprintf("https://%v.%v.app.jexia.com", id, zone),
		token:        Token{},
		tokenRequest: nil,
		// TODO: Add optimisations to default http client
		http: &http.Client{},
	}
	for _, o := range opts {
		o(client)
	}
	return client
}
