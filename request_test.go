package jexiasdkgo

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddHeader(t *testing.T) {
	var request *http.Request
	request = &http.Request{
		Header: http.Header{},
	}
	option := addHeader("title", "value")

	assert.NotEqual(t, request.Header, map[string][]string{})
	option(request)
	assert.Equal(t, request.Header, http.Header(http.Header{"Title": {"value"}}))
}

func TestSetBody(t *testing.T) {
	var request *http.Request
	var target *http.Request

	body := []byte("test body")
	request = &http.Request{}
	target = &http.Request{
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
	}

	option := setBody(body)

	assert.NotEqual(t, target, request)
	option(request)
	assert.Equal(t, target, request)
}

func TestUseAPKMethod(t *testing.T) {
	var request *http.Request
	request = &http.Request{
		Header: http.Header{},
	}
	option := useAPKMethod("keyValue", "secretValue")

	assert.NotEqual(t, request.Header, map[string][]string{})
	option(request)
	assert.Equal(t, request.Header, http.Header(
		http.Header{
			"Method": {"apk"},
			"Key":    {"keyValue"},
			"Secret": {"secretValue"},
		},
	),
	)
}

func TestAddToken(t *testing.T) {
	var request *http.Request
	accessToken := "yourAccessToken"
	request = &http.Request{
		Header: http.Header{},
	}
	option := addToken(accessToken)

	assert.NotEqual(t, request.Header, map[string][]string{})
	option(request)
	assert.Equal(t, request.Header, http.Header(http.Header{"Authorization": {fmt.Sprintf("Bearer %v", accessToken)}}))
}
