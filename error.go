package jexiasdkgo

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Internal signifies the error was send from an internal source
const Internal = "Internal"

// API signifies the error was send from the Jexia API directly
const API = "API"

// Error is the error returned, allows for retry scenarios
type Error struct {
	ID        string `json:"request_id"`
	Message   string `json:"message"`
	Origin    string `json:"origin"`
	Temporary bool   `json:"temporary"`
}

func (e *Error) Error() string {
	temp := ""
	if e.Temporary {
		temp = "temporary"
	}
	return fmt.Sprintf("%v endpoint error: %v (Origin: %v) (ID: %v)", temp, e.Message, e.Origin, e.ID)
}

// checkForAPIError is an internal function wrapper for returning a more useful API error
func checkForAPIError(response *http.Response) error {
	// Success is indicated with 2xx status codes:
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}
	var APIErr []Error
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return &Error{
			ID:        "e001",
			Message:   fmt.Errorf("Unable to read response body: %w", err).Error(),
			Origin:    Internal,
			Temporary: false,
		}
	}
	err = unmarshal(b, &APIErr)
	if err != nil {
		return &Error{
			ID:        "e002",
			Message:   fmt.Errorf("Unable to unmarshal body: %w", err).Error(),
			Origin:    Internal,
			Temporary: false,
		}
	}
	// If there is a matching API error, return this as we can fail absolutely
	if len(APIErr) > 0 {
		return &Error{
			ID:        APIErr[0].ID,
			Message:   APIErr[0].Message,
			Origin:    API,
			Temporary: false,
		}
	}
	// Such unknown error may allow us to retry, it could be due to a network connection drop etc.
	return &Error{
		ID:        "e003",
		Message:   fmt.Errorf("Unknown error, error does not match predefined parameters, presumed failed call").Error(),
		Origin:    Internal,
		Temporary: true,
	}
}

func getNiceError(err error, message string) *Error {
	switch err.(type) {
	case *Error:
		return err.(*Error)
	default:
		return &Error{
			ID:        "e004",
			Message:   fmt.Errorf("%v: %v", message, err).Error(),
			Origin:    Internal,
			Temporary: false,
		}
	}
}
