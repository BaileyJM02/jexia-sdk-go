package jexiasdkgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// APIError is the error returned from a Jexia API endpoint
type APIError struct {
	ID      string `json:"request_id"`
	Message string `json:"message"`
}

// checkForAPIError is an internal function wrapper for returning a more useful API error
func checkForAPIError(response *http.Response) error {
	// Success is indicated with 2xx status codes:
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}
	var error APIError
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Unable to read response body: %w", err)
	}
	json.Unmarshal(b, &error)
	return fmt.Errorf("Endpoint error: %v (ID: %v)", error.Message, error.ID)
}
