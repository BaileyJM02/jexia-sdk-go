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
	)
	assert.Equal(t, "projectID", client.projectID)
	assert.Equal(t, "projectZone", client.projectZone)
	assert.Equal(t, Token{}, client.token)
	assert.Equal(t, nil, client.tokenRequest)
	assert.Equal(t, &http.Client{}, client.http)
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
		SetProjectURL(server.URL),
	)
	assert.Equal(t, client.projectURL, server.URL)
}

func TestNewClientWithAPKToken(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/auth")
		assert.Equal(t, req.Method, http.MethodPost)

		actual := APKTokenRequest{}
		target := APKTokenRequest{
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
		SetProjectURL(server.URL),
	)
	client.UseAPKToken("APIKey", "APISecret")
	assert.Equal(t, "yourAccessToken", client.token.Access)
	assert.Equal(t, "yourRefreshToken", client.token.Refresh)
	assert.Equal(t, 118*time.Minute, client.token.Lifetime)
}

func TestRefreshTokenWithAPKToken(t *testing.T) {
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
		token:       token,
		tokenRequest: APKTokenRequest{
			Method: "apk",
			Key:    "APIKey",
			Secret: "APISecret",
		},
		http: &http.Client{},
	}

	client.RefreshToken()
	assert.Equal(t, "yourCurrentAccessToken", client.token.Access)
	assert.Equal(t, "yourNewRefreshToken", client.token.Refresh)
}

func TestSetTokenLifetimeWithAPKToken(t *testing.T) {
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
		SetProjectURL(server.URL),
	)
	client.UseAPKToken("APIKey", "APISecret")
	assert.Equal(t, 118*time.Minute, client.token.Lifetime)

	client.SetTokenLifetime(60 * time.Minute)
	assert.Equal(t, 60*time.Minute, client.token.Lifetime)

	assert.Equal(t, "yourAccessToken", client.token.Access)
	assert.Equal(t, "yourRefreshToken", client.token.Refresh)
}

func TestNewClientWithUMSToken(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, req.URL.String(), "/auth")
		assert.Equal(t, req.Method, http.MethodPost)

		actual := UMSTokenRequest{}
		target := UMSTokenRequest{
			Method:   "ums",
			Email:    "email",
			Password: "password",
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
		SetProjectURL(server.URL),
	)
	client.UseUMSToken("email", "password")
	client.AutoRefreshToken()

	assert.Equal(t, "yourAccessToken", client.token.Access)
	assert.Equal(t, "yourRefreshToken", client.token.Refresh)
	assert.Equal(t, 118*time.Minute, client.token.Lifetime)
}

func TestRefreshTokenWithUMSToken(t *testing.T) {
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
		token:       token,
		tokenRequest: UMSTokenRequest{
			Method:   "ums",
			Email:    "APIKey",
			Password: "APISecret",
		},
		http: &http.Client{},
	}

	client.RefreshToken()
	assert.Equal(t, "yourCurrentAccessToken", client.token.Access)
	assert.Equal(t, "yourNewRefreshToken", client.token.Refresh)
}

func TestSetTokenLifetimeWithUMSToken(t *testing.T) {
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
		SetProjectURL(server.URL),
	)
	client.UseUMSToken("email", "password")
	assert.Equal(t, 118*time.Minute, client.token.Lifetime)

	client.SetTokenLifetime(60 * time.Minute)
	assert.Equal(t, 60*time.Minute, client.token.Lifetime)

	assert.Equal(t, "yourAccessToken", client.token.Access)
	assert.Equal(t, "yourRefreshToken", client.token.Refresh)
}

func TestAutoRefreshToken(t *testing.T) {
	token := Token{
		Access:   "yourCurrentAccessToken",
		Refresh:  "yourCurrentRefreshToken",
		Lifetime: time.Millisecond,
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
	defer	server.Close()

	client := &Client{
		projectID:   "projectID",
		projectZone: "projectZone",
		projectURL:  server.URL,
		token:       token,
		tokenRequest: UMSTokenRequest{
			Method:   "ums",
			Email:    "APIKey",
			Password: "APISecret",
		},
		http: &http.Client{},
		abortRefresh: make(chan bool),
	}

	client.AutoRefreshToken()
	// The shortest amount of time we can wait before a race condition is likely to appear due to the short periods of time
	// The fastest so we don't have hanging tests
	time.Sleep(3950 * time.Microsecond)
	// stop the refresh loop externally
	close(client.abortRefresh)
	assert.Equal(t, "yourCurrentAccessToken", client.GetToken().Access)
	assert.Equal(t, "yourNewRefreshToken", client.GetToken().Refresh)
}
func TestNewRefreshCycle(t *testing.T) {
	token := Token{
		Access:   "yourCurrentAccessToken",
		Refresh:  "yourCurrentRefreshToken",
		Lifetime: time.Millisecond,
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
	defer	server.Close()

	client := &Client{
		projectID:   "projectID",
		projectZone: "projectZone",
		projectURL:  server.URL,
		token:       token,
		tokenRequest: UMSTokenRequest{
			Method:   "ums",
			Email:    "APIKey",
			Password: "APISecret",
		},
		http: &http.Client{},
		abortRefresh: make(chan bool),
	}

	client.newRefreshCycle()
	// The shortest amount of time we can wait before a race condition is likely to 
	// appear due to the short periods of time but the fastest so we don't have hanging tests
	time.Sleep(3950 * time.Microsecond)
	// stop the refresh loop externally
	close(client.abortRefresh)

	assert.Equal(t, "yourCurrentAccessToken", client.GetToken().Access)
	assert.Equal(t, "yourNewRefreshToken", client.GetToken().Refresh)
}

