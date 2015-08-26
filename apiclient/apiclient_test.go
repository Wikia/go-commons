package apiclient

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	BaseURL     = "http://wikia.com"
	Endpoint    = "test"
	EndpointURL = BaseURL + "/" + Endpoint
	BadEndpoint = ":"
	HeaderTest  = "http://headers.jsontest.com/"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(BaseURL)

	assert.Nil(t, err, "NewClient creation error")

	assert.Equal(t, BaseURL, client.BaseURL.String(), "NewClient BaseURL")
}

func TestNewRequest(t *testing.T) {
	client, _ := NewClient(BaseURL)

	request, err := client.NewRequest("GET", Endpoint, nil)

	assert.Nil(t, err, "NewRequest creation error")

	assert.Equal(t, EndpointURL, request.URL.String(), "NewRequest endpoint URL")

	_, err = client.NewRequest("GET", BadEndpoint, nil)

	assert.NotNil(t, err, "NewRequest bad URL no error")
}

func TestCallWithHeaders(t *testing.T) {
	client, _ := NewClient(HeaderTest)

	headers := map[string]string{
		"custom-header":  "foo",
		"Another-Header": "bar",
	}

	response, err := client.Call("GET", "", nil, headers)

	assert.NoError(t, err, "Error getting Header response")

	var f interface{}
	err = json.NewDecoder(response.Body).Decode(&f)
	assert.NoError(t, err, fmt.Sprintf("Error deserializing JSON: %#v", response.Body))

	data := f.(map[string]interface{})
	assert.NotNil(t, data["Custom-Header"], "Custom-Header is missing")
	assert.Equal(t, "foo", data["Custom-Header"].(string), "Custom-Header is invalid")
	assert.NotNil(t, data["Another-Header"], "Another-Header is missing")
	assert.Equal(t, "bar", data["Another-Header"].(string), "Another-Header is invalid")
}

func TestCallWithoutHeaders(t *testing.T) {
	client, _ := NewClient(HeaderTest)

	response, err := client.Call("GET", "", nil, nil)

	assert.NoError(t, err, "Error getting Header response")

	var f interface{}
	err = json.NewDecoder(response.Body).Decode(&f)
	assert.NoError(t, err, fmt.Sprintf("Error deserializing JSON: %#v", response.Body))

	data := f.(map[string]interface{})
	assert.Equal(t, 3, len(data), "Incorrect number of headers sent")

	for k, v := range data {
		fmt.Printf("%v: %v\n", k, v.(string))
	}
	assert.NotNil(t, data["Content-Type"], "Content-Type is missing")
	assert.Equal(t, "application/x-www-form-urlencoded", data["Content-Type"].(string), "Content-Type is invalid")

	assert.NotNil(t, data["User-Agent"], "User-Agent is missing")
	assert.Equal(t, "Go-http-client/1.1", data["User-Agent"].(string), "User-Agent is invalid")

	assert.NotNil(t, data["Host"], "Host is missing")
	assert.Equal(t, "headers.jsontest.com", data["Host"].(string), "Host is invalid")
}
