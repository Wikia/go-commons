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
}

type phalanxClientImpl struct {
	apiClient *apiclient.Client
}

type checkResult struct {
	contentType string
	content     string
	ok          bool
	err         error
}

func NewClient(baseURL string) (PhalanxClient, error) {
	apiClient, err := apiclient.NewClient(baseURL)
	if err != nil {
		return nil, err
	}
	client := &phalanxClientImpl{apiClient: apiClient}
	return client, nil
}

func (client *phalanxClientImpl) CheckName(name string) (bool, error) {
	return client.check(checkTypeName, name)
}

func (client *phalanxClientImpl) CheckNames(names ...string) (bool, error) {
	return client.checkMultiple(checkTypeName, names...)
}

func (client *phalanxClientImpl) CheckNamesConcurrent(names ...string) (bool, error) {
	return client.checkMultipleConcurrent(checkTypeName, names...)
}

func (client *phalanxClientImpl) CheckEmail(email string) (bool, error) {
	return client.check(checkTypeEmail, email)
}

func (client *phalanxClientImpl) CheckEmails(names ...string) (bool, error) {
	return client.checkMultiple(checkTypeEmail, names...)
}

func (client *phalanxClientImpl) CheckEmailsConcurrent(names ...string) (bool, error) {
	return client.checkMultipleConcurrent(checkTypeEmail, names...)
}

func (client *phalanxClientImpl) check(checkType, content string) (bool, error) {
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

func (client *phalanxClientImpl) checkMultiple(checkType string, content ...string) (bool, error) {
	for _, c := range content {
		ok, err := client.check(checkType, c)

		if err != nil || !ok {
			return ok, err
		}
	}

	return true, nil
}

func (client *phalanxClientImpl) checkMultipleConcurrent(checkType string, content ...string) (bool, error) {
	channels := make([]<-chan *checkResult, 0)

	for _, c := range content {
		channels = append(channels, client.checkConcurrent(checkType, c))
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

func (client *phalanxClientImpl) checkConcurrent(contentType string, content string) <-chan *checkResult {
	out := make(chan *checkResult)

	go func() {
		defer close(out)
		ok, err := client.check(contentType, content)
		out <- &checkResult{
			ok:          ok,
			err:         err,
			content:     content,
			contentType: contentType,
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
