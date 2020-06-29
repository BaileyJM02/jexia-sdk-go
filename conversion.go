package jexiasdkgo

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// marshal is an internal function wrapper for marshalling json payloads, may be more intricate in the future
func marshal(payload interface{}) ([]byte, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// unmarshal is an internal function wrapper for unmarshalling json payloads from a ReadCloser, may be more intricate in the future
func unmarshal(b []byte, target interface{}) error {
	json.Unmarshal(b, &target)
	return nil
}

// unmarshal is an internal function wrapper for unmarshalling json payloads from a ReadCloser, may be more intricate in the future
func read(body io.Reader) ([]byte, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	return b, err
}
