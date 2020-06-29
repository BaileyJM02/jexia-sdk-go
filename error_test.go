package jexiasdkgo

import (
	"bytes"
	"net/http"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckErrorOnFailedCall(t *testing.T) {
	var response *http.Response
	var target APIError

	target = APIError{
		ID:      "some-really-long-id",
		Message: "A really useful message",
	}

	body, err := marshal(target)
	if err != nil {
		assert.Error(t, err)
	}

	response = &http.Response{
		StatusCode:    400,
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
	}

	err = checkForAPIError(response)
	assert.Contains(t, err.Error(), target.ID)
	assert.Contains(t, err.Error(), target.Message)
}

func TestCheckErrorOnSuccessfulCall(t *testing.T) {
	var response *http.Response
	var target APIError
	
	target = APIError{
		ID:      "some-really-long-id",
		Message: "A really useful message",
	}

	body, err := marshal(target)
	if err != nil {
		assert.Error(t, err)
	}

	response = &http.Response{
		StatusCode:    200,
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
	}

	err = checkForAPIError(response)
	assert.Equal(t, err, nil)
}
