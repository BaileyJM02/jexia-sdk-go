package jexiasdkgo

import (
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
		// Send response to be tested
		rw.Write([]byte(`{"access_token":"yourAccessToken","refresh_token":"yourRefreshToken"}`))
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
}
