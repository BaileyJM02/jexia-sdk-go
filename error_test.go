package jexiasdkgo

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckErrorOnFailedCall(t *testing.T) {
	var response *http.Response
	var target []APIError

	target = []APIError{APIError{
		ID:      "some-really-long-id",
		Message: "A really useful message",
	}}

	body, err := marshal(target)
	if err != nil {
		assert.Error(t, err)
	}

	response = &http.Response{
		StatusCode: 400,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}

	err = checkForAPIError(response)
	assert.Contains(t, err.Error(), target[0].ID)
	assert.Contains(t, err.Error(), target[0].Message)
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
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader(body)),
	}

	err = checkForAPIError(response)
	assert.Equal(t, err, nil)
}
