package apiclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	BaseURL     = "http://wikia.com"
	Endpoint    = "test"
	EndpointURL = BaseURL + "/" + Endpoint
	BadEndpoint = ":"
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
