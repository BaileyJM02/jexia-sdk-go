package jexiasdkgo

import (
	"fmt"
	"sync"
)

// Dataset a struct containing the name of the dataset and the client memory pointer
// As the client is a memory pointer, any changes made to the client will be reflected within the dataset, therefore token refreshes will still work
type Dataset struct {
	Name   string
	Client *Client
	mux    sync.Mutex
}

// GetDataset returns a dataset instance that can be used to perform actions against
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

// Insert allows you to add sets of data to the dataset
// If you are inputing a single interface you must wrap it in an array: []interface{}{yourStruct}
func (d *Dataset) Insert(dataArray []interface{}) ([]interface{}, error) {
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
	return result, nil
}

// Select allows you to select the data from a dataset
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
