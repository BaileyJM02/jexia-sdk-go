package jexiasdkgo

import (
	"fmt"
	"sync"
)

// Dataset
type Dataset struct {
	Name   string
	Client *Client
	mux    sync.Mutex
}

// GetDataset
func (c *Client) GetDataset(name string) *Dataset {
	return &Dataset{
		Name:   name,
		Client: c,
	}
}

// GetName fetches the dataset name
func (d *Dataset) GetName() string {
	d.mux.Lock()
	name := d.Name
	d.mux.Unlock()
	return name
}

// GetClient fetches the current client
func (d *Dataset) GetClient() *Client {
	d.mux.Lock()
	client := d.Client
	d.mux.Unlock()
	return client
}

// Insert
func (d *Dataset) Insert(dataArray []interface{}) (interface{}, error) {
	var result []interface{}
	payload, err := marshal(dataArray)
	if err != nil {
		return nil, err
	}
	err = d.GetClient().post(
		fmt.Sprintf("%v/ds/%v", d.GetClient().projectURL, d.GetName()),
		&result,
		addToken(d.GetClient().GetToken().Access),
		setBody(payload),
	)
	if err != nil {
		return nil, err
	}
	return result[0], nil
}

// Select
func (d *Dataset) Select() ([]interface{}, error) {
	var result []interface{}
	err := d.GetClient().get(
		fmt.Sprintf("%v/ds/%v", d.GetClient().projectURL, d.GetName()),
		&result,
		addToken(d.GetClient().GetToken().Access),
	)
	if err != nil {
		return nil, err
	}
	return result, nil
}
