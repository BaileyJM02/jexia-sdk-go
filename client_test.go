package jexiasdkgo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetHTTPClient(t *testing.T) {
	var client *Client
	var httpClient *http.Client

	client = NewClient(
		"projectID",
		"projectZone",
		"APIKey",
		"APISecret",
	)
	httpClient = &http.Client{nil, nil, nil, 2 * time.Microsecond}

	option := SetHTTPClient(httpClient)

	assert.NotEqual(t, httpClient, client.http)
	option(client)
	assert.Equal(t, httpClient, client.http)
}

func TestSetProjectURL(t *testing.T) {
	var client *Client

	client = NewClient(
		"projectID",
		"projectZone",
		"APIKey",
		"APISecret",
	)

	option := SetProjectURL("testURL")

	assert.NotEqual(t, "testURL", client.projectURL)
	option(client)
	assert.Equal(t, "testURL", client.projectURL)
}

func TestNewClient(t *testing.T) {
	var client *Client
	client = NewClient(
		"projectID",
		"projectZone",
		"APIKey",
		"APISecret",
	)
	assert.Equal(t, "projectID", client.projectID)
	assert.Equal(t, "projectZone", client.projectZone)
	assert.Equal(t, "APIKey", client.apiKey)
	assert.Equal(t, "APISecret", client.apiSecret)
	assert.Equal(t, Token{}, client.token)
}

func TestSetNewClientProjectUrl(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Fail(t, "Server should not be called at NewClient at any point.")
	}))
	// Close the server when test finishes
	defer server.Close()

	var client *Client
	client = NewClient(
		"projectID",
		"projectZone",
		"APIKey",
		"APISecret",
		SetProjectURL(server.URL),
	)
	assert.Equal(t, client.projectURL, server.URL)
}

func TestNewClientWithToken(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/auth")
		assert.Equal(t, req.Method, http.MethodPost)

		actual := APITokenRequest{}
		target := APITokenRequest{
			Method: "apk",
			Key:    "APIKey",
			Secret: "APISecret",
		}
		b, err := read(req.Body)
		if err != nil {
			assert.Error(t, err)
		}
		err = unmarshal(b, &actual)
		if err != nil {
			assert.Error(t, err)
		}
		assert.Equal(t, target, actual)

		payload, _ := marshal(Token{
			Access:  "yourAccessToken",
			Refresh: "yourRefreshToken",
		})
		// Send response to be tested
		rw.Write(payload)
	}))
	// Close the server when test finishes
	defer server.Close()

	var client *Client
	client = NewClient(
		"projectID",
		"projectZone",
		"APIKey",
		"APISecret",
		SetProjectURL(server.URL),
	)
	client.GetToken()
	assert.Equal(t, "yourAccessToken", client.token.Access)
	assert.Equal(t, "yourRefreshToken", client.token.Refresh)
	assert.Equal(t, 118*time.Minute, client.token.Lifetime)
}

func TestSetTokenLifetime(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/auth")
		assert.Equal(t, req.Method, http.MethodPost)

		payload, _ := marshal(Token{
			Access:  "yourAccessToken",
			Refresh: "yourRefreshToken",
		})
		// Send response to be tested
		rw.Write(payload)
	}))
	// Close the server when test finishes
	defer server.Close()

	var client *Client
	client = NewClient(
		"projectID",
		"projectZone",
		"APIKey",
		"APISecret",
		SetProjectURL(server.URL),
	)
	client.GetToken()
	assert.Equal(t, 118*time.Minute, client.token.Lifetime)

	client.SetTokenLifetime(60 * time.Minute)
	assert.Equal(t, 60*time.Minute, client.token.Lifetime)

	assert.Equal(t, "yourAccessToken", client.token.Access)
	assert.Equal(t, "yourRefreshToken", client.token.Refresh)
}

func TestRefreshToken(t *testing.T) {
	token := Token{
		Access:  "yourCurrentAccessToken",
		Refresh: "yourCurrentRefreshToken",
	}
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/auth/refresh", req.URL.String())
		headers := req.Header
		assert.Equal(t, 1, len(headers["Authorization"]))
		assert.Equal(t, fmt.Sprintf("Bearer %v", token.Access), headers["Authorization"][0])
		// Send response to be tested
		assert.Equal(t, req.Method, http.MethodPost)

		payload, _ := marshal(Token{
			Refresh: "yourNewRefreshToken",
		})
		// Send response to be tested
		rw.Write(payload)
	}))
	// Close the server when test finishes
	defer server.Close()

	client := &Client{
		projectID:   "projectID",
		projectZone: "projectZone",
		projectURL:  server.URL,
		apiKey:      "APIKey",
		apiSecret:   "APISecret",
		token:       token,
		http:        &http.Client{},
	}

	client.RefreshToken()
	assert.Equal(t, "yourCurrentAccessToken", client.token.Access)
	assert.Equal(t, "yourNewRefreshToken", client.token.Refresh)
}
