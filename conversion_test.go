package jexiasdkgo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshal(t *testing.T) {
	type test struct {
		value string
	}

	input := test{
		"testing",
	}

	actualPure, err := json.Marshal(input)
	if err != nil {
		assert.Error(t, err)
	}

	actualMarshal, err := marshal(input)
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, actualPure, actualMarshal)
}

func TestUnmarshal(t *testing.T) {
	type test struct {
		value string
	}
	var actualPure test
	var actualMarshal test

	input := []byte("value: string")

	err := json.Unmarshal(input, &actualPure)
	if err != nil {
		assert.Error(t, err)
	}

	err = unmarshal(input, &actualMarshal)
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, actualPure, actualMarshal)
}

func TestRead(t *testing.T) {
	// Generate two as the first call will empty the array for the second
	bodyTarget := ioutil.NopCloser(bytes.NewReader([]byte("test string")))
	bodyActual := ioutil.NopCloser(bytes.NewReader([]byte("test string")))

	target, err := ioutil.ReadAll(bodyTarget)
	if err != nil {
		assert.Error(t, err)
	}

	actual, err := read(bodyActual)
	if err != nil {
		assert.Error(t, err)
	}

	assert.Equal(t, target, actual)
}
