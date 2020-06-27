package jexiasdkgo

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	var client *Client
	client = NewClient(
		os.Getenv("PROJECT_ID"),
		os.Getenv("PROJECT_ZONE"),
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
	)
	assert.Equal(t, client.projectID, os.Getenv("PROJECT_ID"))
	assert.Equal(t, client.projectZone, os.Getenv("PROJECT_ZONE"))
	assert.Equal(t, client.apiKey, os.Getenv("API_KEY"))
	assert.Equal(t, client.apiSecret, os.Getenv("API_SECRET"))
	assert.Equal(t, client.token, Token{})
}

func TestSetNewClientProjectUrl(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Server should not be called with NewClient at any point.
		assert.True(t, false)
	}))
	// Close the server when test finishes
	defer server.Close()

	var client *Client
	client = NewClient(
		os.Getenv("PROJECT_ID"),
		os.Getenv("PROJECT_ZONE"),
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
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
		os.Getenv("PROJECT_ID"),
		os.Getenv("PROJECT_ZONE"),
		os.Getenv("API_KEY"),
		os.Getenv("API_SECRET"),
		SetProjectURL(server.URL),
	)
	client.GetToken()
	assert.Equal(t, "yourAccessToken", client.token.Access)
	assert.Equal(t, "yourRefreshToken", client.token.Refresh)
}
