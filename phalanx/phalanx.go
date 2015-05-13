package phalanx

import (
	"net/url"

	"github.com/Wikia/go-commons/apiclient"
)

const (
	contentKey     = "content"
	typeKey        = "type"
	checkEndpoint  = "check"
	checkOk        = "ok\n"
	checkTypeName  = "user"
	checkTypeEmail = "email"
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

func (client *Client) CheckName(name string) (bool, error) {
	return client.Check(checkTypeName, name)
}

func (client *Client) CheckEmail(email string) (bool, error) {
	return client.Check(checkTypeEmail, email)
}

func (client *Client) Check(checkType, content string) (bool, error) {

	data := url.Values{}
	data.Add(typeKey, checkType)
	data.Add(contentKey, content)

	resp, err := client.apiClient.Call("POST", checkEndpoint, data)
	if err != nil {
		return false, err
	}

	resBody, err := client.apiClient.GetBody(resp)
	if err != nil {
		return false, err
	}

	if string(resBody) == checkOk {
		return true, nil
	}

	return false, nil
}
