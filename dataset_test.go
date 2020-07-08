package jexiasdkgo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetNameFromDataset(t *testing.T) {
	var dataset *Dataset
	dataset = &Dataset{
		Name: "name",
	}
	assert.Equal(t, "name", dataset.GetName())
}

func TestGetClientFromDataset(t *testing.T) {
	var dataset *Dataset
	var client *Client
	client = &Client{
		projectID: "testing",
	}
	dataset = &Dataset{
		Client: client,
	}
	assert.Equal(t, client, dataset.GetClient())
}

func TestClientPassedAsMemoryPointerToDataset(t *testing.T) {
	var client *Client
	client = NewClient(
		"projectID",
		"projectZone",
	)

	client.SetToken(Token{
		Access:  "yourCurrentAccessToken",
		Refresh: "yourCurrentRefreshToken",
	})
	assert.Equal(t, "yourCurrentAccessToken", client.GetToken().Access)
	assert.Equal(t, "yourCurrentRefreshToken", client.GetToken().Refresh)

	var dataset *Dataset
	dataset = client.GetDataset("datasetName")
	assert.Equal(t, "datasetName", dataset.GetName())
	assert.Equal(t, client, dataset.GetClient())

	client.SetToken(Token{
		Access:  "yourNewAccessToken",
		Refresh: "yourNewRefreshToken",
	})
	assert.Equal(t, "yourNewAccessToken", client.GetToken().Access)
	assert.Equal(t, "yourNewRefreshToken", client.GetToken().Refresh)

	assert.Equal(t, "datasetName", dataset.GetName())
	assert.Equal(t, client, dataset.GetClient())
	// Directly check to be 100% certain
	assert.Equal(t, "yourNewAccessToken", dataset.GetClient().GetToken().Access)
}

func TestDatasetSelect(t *testing.T) {
	var token Token
	token = Token{
		Access: "yourCurrentAccessToken",
	}
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/ds/test", req.URL.String())
		headers := req.Header
		assert.Equal(t, 1, len(headers["Authorization"]))
		assert.Equal(t, fmt.Sprintf("Bearer %v", token.Access), headers["Authorization"][0])
		// Send response to be tested
		assert.Equal(t, http.MethodGet, req.Method)

		payload := ([]byte(`[{"id":"test","created_at":"2020-07-08T16:08:50.304789Z","updated_at":"2020-07-08T16:08:50.304789Z","@type":"some-type","@name":"tabletop"},{"id":"5d7b907c-06bd-41e3-a113-addc230635e1","created_at":"2020-07-08T16:09:30.143354Z","updated_at":"2020-07-08T16:09:30.143354Z","@type":"some-type","@name":"tabletop"},{"id":"d8243cb2-6b87-4d19-8b36-ad51c52103db","created_at":"2020-07-08T16:10:32.089546Z","updated_at":"2020-07-08T16:10:32.089546Z","@type":"some-typ42e","@name":"table324top"}]`))
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
	client.SetToken(token)
	dataset := client.GetDataset("test")
	data, err := dataset.Select()
	if err != nil {
		assert.Error(t, err)
	}
	var expected interface{}
	type expectedStruct struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Type      string    `json:"@type"`
		Name      string    `json:"@name"`
	}

	var actualStruct expectedStruct
	expStruct := expectedStruct{
		ID:        "test",
		CreatedAt: time.Date(2020, 07, 8, 16, 8, 50, 304789000, time.UTC),
		UpdatedAt: time.Date(2020, 07, 8, 16, 8, 50, 304789000, time.UTC),
		Type:      "some-type",
		Name:      "tabletop",
	}

	err = unmarshal([]byte(`{"id":"test","created_at":"2020-07-08T16:08:50.304789Z","updated_at":"2020-07-08T16:08:50.304789Z","@type":"some-type","@name":"tabletop"}`), &expected)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, expected, data[0])

	err = unmarshal([]byte(`{"id":"test","created_at":"2020-07-08T16:08:50.304789Z","updated_at":"2020-07-08T16:08:50.304789Z","@type":"some-type","@name":"tabletop"}`), &actualStruct)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, expStruct, actualStruct)
}

func TestDatasetInsert(t *testing.T) {
	var token Token
	token = Token{
		Access: "yourCurrentAccessToken",
	}
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		assert.Equal(t, "/ds/test", req.URL.String())
		headers := req.Header
		assert.Equal(t, 1, len(headers["Authorization"]))
		assert.Equal(t, fmt.Sprintf("Bearer %v", token.Access), headers["Authorization"][0])
		// Send response to be tested
		assert.Equal(t, http.MethodPost, req.Method)

		payload := ([]byte(`[{"id":"test","created_at":"2020-07-08T16:08:50.304789Z","updated_at":"2020-07-08T16:08:50.304789Z","@type":"some-type","@name":"tabletop"},{"id":"5d7b907c-06bd-41e3-a113-addc230635e1","created_at":"2020-07-08T16:09:30.143354Z","updated_at":"2020-07-08T16:09:30.143354Z","@type":"some-type","@name":"tabletop"},{"id":"d8243cb2-6b87-4d19-8b36-ad51c52103db","created_at":"2020-07-08T16:10:32.089546Z","updated_at":"2020-07-08T16:10:32.089546Z","@type":"some-typ42e","@name":"table324top"}]`))
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
	client.SetToken(token)
	dataset := client.GetDataset("test")

	var expected interface{}
	type expectedStruct struct {
		ID        string    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Type      string    `json:"@type"`
		Name      string    `json:"@name"`
	}

	var actualStruct expectedStruct
	expStruct := expectedStruct{
		ID:        "test",
		CreatedAt: time.Date(2020, 07, 8, 16, 8, 50, 304789000, time.UTC),
		UpdatedAt: time.Date(2020, 07, 8, 16, 8, 50, 304789000, time.UTC),
		Type:      "some-type",
		Name:      "tabletop",
	}

	data, err := dataset.Insert([]interface{}{expStruct})
	if err != nil {
		assert.Error(t, err)
	}

	err = unmarshal([]byte(`{"id":"test","created_at":"2020-07-08T16:08:50.304789Z","updated_at":"2020-07-08T16:08:50.304789Z","@type":"some-type","@name":"tabletop"}`), &expected)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, expected, data[0])

	err = unmarshal([]byte(`{"id":"test","created_at":"2020-07-08T16:08:50.304789Z","updated_at":"2020-07-08T16:08:50.304789Z","@type":"some-type","@name":"tabletop"}`), &actualStruct)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, expStruct, actualStruct)
}
