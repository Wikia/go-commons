package phalanx

import (
	"net/url"

	"github.com/Wikia/go-commons/apiclient"
	"sync"
)

const (
	contentKey     = "content"
	typeKey        = "type"
	checkEndpoint  = "check"
	checkOk        = "ok\n"
	checkTypeName  = "user"
	checkTypeEmail = "email"
)

type PhalanxClient interface {
	CheckName(name string) (bool, error)
	CheckNames(names ...string) (bool, error)
	CheckNamesConcurrent(names ...string) (bool, error)
	CheckEmail(email string) (bool, error)
	CheckEmails(names ...string) (bool, error)
	CheckEmailsConcurrent(names ...string) (bool, error)
	Check(checkType, content string) (bool, error)
	CheckMultiple(checkType string, content ...string) (bool, error)
	CheckMultipleConcurrent(checkType string, content ...string) (bool, error)
}

type client struct {
	apiClient *apiclient.Client
}

type checkResult struct {
	contentType string
	ok          bool
	err         error
}

func NewClient(baseURL string) (PhalanxClient, error) {
	apiClient, err := apiclient.NewClient(baseURL)
	if err != nil {
		return nil, err
	}
	client := &client{apiClient: apiClient}
	return client, nil
}

func (client *client) CheckName(name string) (bool, error) {
	return client.Check(checkTypeName, name)
}

func (client *client) CheckNames(names ...string) (bool, error) {
	return client.CheckMultiple(checkTypeName, names...)
}

func (client *client) CheckNamesConcurrent(names ...string) (bool, error) {
	return client.CheckMultipleConcurrent(checkTypeName, names...)
}

func (client *client) CheckEmail(email string) (bool, error) {
	return client.Check(checkTypeEmail, email)
}

func (client *client) CheckEmails(names ...string) (bool, error) {
	return client.CheckMultiple(checkTypeEmail, names...)
}

func (client *client) CheckEmailsConcurrent(names ...string) (bool, error) {
	return client.CheckMultipleConcurrent(checkTypeEmail, names...)
}

func (client *client) Check(checkType, content string) (bool, error) {
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

func (client *client) CheckMultiple(checkType string, content ...string) (bool, error) {
	for _, c := range content {
		ok, err := client.Check(checkType, c)

		if err != nil || !ok {
			return ok, err
		}
	}

	return true, nil
}

func (client *client) CheckMultipleConcurrent(checkType string, content ...string) (bool, error) {
	channels := make([]<-chan *checkResult)

	for _, c := range content {
		channels = append(channels, client.checkContentConcurrent(checkType, c))
	}

	for check := range waitForConcurrentRequests(channels...) {
		if check.err != nil {
			return false, check.err
		} else if !check.ok {
			return false, nil
		}
	}

	return true, nil
}

func (client *client) checkContentConcurrent(checkType string, content string) <-chan *checkResult {
	out := make(chan *checkResult)

	go func() {
		defer close(out)
		ok, err := client.Check(checkType, content)
		out <- &checkResult{
			ok:      ok,
			err:     err,
			content: content,
		}
	}()

	return out
}

func waitForConcurrentRequests(channels ...<-chan *checkResult) <-chan *checkResult {
	var waitGroup sync.WaitGroup
	out := make(chan *checkResult)

	copyToOutput := func(channel <-chan *checkResult) {
		defer waitGroup.Done()
		for check := range channel {
			out <- check
		}
	}

	waitGroup.Add(len(channels))
	for _, channel := range channels {
		go copyToOutput(channel)
	}

	go func() {
		defer close(out)
		waitGroup.Wait()
	}()

	return out
}
