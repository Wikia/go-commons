package phalanx

import (
	"github.com/Wikia/go-commons/apiclient"
	"net/url"
)

const (
	contentKey    = "content"
	typeKey       = "type"
	typeUser      = "user"
	checkEndpoint = "check"
	checkOk       = "ok\n"
)

type Client struct {
	apiClient *apiclient.Client
}

func NewClient(baseURL string) (*Client, error) {
	apiClient, err := apiclient.NewClient(baseURL)
	if err != nil {
		return nil, err
	}
	client := &Client{apiClient: apiClient}
	return client, nil
}

func (client *Client) Check(name string) (bool, error) {

	data := url.Values{}
	data.Add(typeKey, typeUser)
	data.Add(contentKey, name)

	resBody, err := client.apiClient.Call("POST", checkEndpoint, data)
	if err != nil {
		return false, err
	}

	if string(resBody) == checkOk {
		return true, nil
	}

	return false, nil
}
