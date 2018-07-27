package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	BaseURL     = "http://wikia.com"
	ProxyURL    = "http://dev.icache.service.sjc-dev.consul:80"
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

func TestNewRequestWithProxy(t *testing.T) {
	client, _ := NewClientWithProxy(BaseURL, ProxyURL)

	request, err := client.NewRequest("GET", Endpoint, nil)

	assert.Nil(t, err, "NewRequest creation error")

	assert.Equal(t, EndpointURL, request.URL.String(), "NewRequest endpoint URL")
}

func TestNewRequestWithProxyBadEndpoint(t *testing.T) {
	client, _ := NewClientWithProxy(BaseURL, ProxyURL)

	_, err := client.NewRequest("GET", BadEndpoint, nil)

	assert.NotNil(t, err, "NewRequest bad URL no error")
}

func TestIfHttpClientsAreDifferent(t *testing.T) {
	client1, _ := NewClientWithProxy(BaseURL, ProxyURL)
	client2, _ := NewClientWithProxy(BaseURL, ProxyURL)

	assert.NotEqual(t, client1.httpClient, client2.httpClient)
}

func TestCallWithHeaders(t *testing.T) {
	client, _ := NewClient(HeaderTest)

	headers := http.Header{
		"custom-header":  []string{"foo"},
		"Another-Header": []string{"bar"},
	}

	response, err := client.Call("GET", "", nil, headers)

	assert.NoError(t, err, "Error getting Header response")

	var f interface{}
	err = json.NewDecoder(response.Body).Decode(&f)
	assert.NoError(t, err, fmt.Sprintf("Error deserializing JSON: %#v", response.Body))

	data := f.(map[string]interface{})
	assert.NotNil(t, data["Custom-Header"], "custom-header is missing")
	assert.Equal(t, "foo", data["Custom-Header"].(string), "custom-header is invalid")
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
	assert.True(t, len(data) >= 3, "Incorrect number of headers sent")

	assert.NotNil(t, data["User-Agent"], "User-Agent is missing")
	assert.Equal(t, "Go-http-client/1.1", data["User-Agent"].(string), "User-Agent is invalid")

	assert.NotNil(t, data["Host"], "Host is missing")
	assert.Equal(t, "headers.jsontest.com", data["Host"].(string), "Host is invalid")

	assert.NotNil(t, data["X-Cloud-Trace-Context"], "X-Cloud-Trace-Context is missing")
}
