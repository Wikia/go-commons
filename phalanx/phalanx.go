package phalanx

import (
	"encoding/json"
	"github.com/Wikia/go-commons/apiclient"
	"net/url"
)

const (
	contentKey     = "content"
	typeKey        = "type"
	checkEndpoint  = "check"
	matchEndpoint  = "match"
	checkOk        = "ok\n"
	checkTypeName  = "user"
	checkTypeEmail = "email"
)

type MatchRecord struct {
	Regex         bool   `json:"regex"`
	Expires       string `json:"expires"`
	Text          string `json:"text"`
	Reason        string `json:"reason"`
	Exact         bool   `json:"exact"`
	CaseSensitive bool   `json:"caseSensative"`
	Id            int    `json:"id"`
	Language      string `json:"language"`
	AuthorId      int    `json:"authorId"`
	Type          int    `json:"type"`
}

type PhalanxClient interface {
	CheckName(name string) (bool, error)
	CheckEmail(email string) (bool, error)
	Check(checkType, content string) (bool, error)
	Match(matchType, content string) ([]MatchRecord, error)
}

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

	resp, err := client.apiClient.Call("POST", checkEndpoint, data, map[string]string{})
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

func (client *Client) Match(checkType, content string) ([]MatchRecord, error) {
	response := make([]MatchRecord, 0)

	data := url.Values{}
	data.Add(typeKey, checkType)
	data.Add(contentKey, content)

	resp, err := client.apiClient.Call("POST", matchEndpoint, data, map[string]string{})
	if err != nil {
		return nil, err
	}

	resBody, err := client.apiClient.GetBody(resp)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
