package phalanx

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/Wikia/go-commons/apiclient"
	"github.com/Wikia/go-commons/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"golang.org/x/net/context"
)

const (
	contentKey     = "content"
	typeKey        = "type"
	checkEndpoint  = "check"
	matchEndpoint  = "match"
	checkOk        = "ok\n"
	CheckTypeName  = "user"
	CheckTypeEmail = "email"

	retriesLimit = 3
	retrySleep   = 20
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
	CheckName(ctx context.Context, name string) (bool, error)
	CheckEmail(ctx context.Context, email string) (bool, error)
	Check(ctx context.Context, checkType, content string) (bool, error)
	Match(ctx context.Context, matchType, content string) ([]MatchRecord, error)
}

type Client struct {
	apiClient apiclient.ApiClient
}

func NewClient(baseURL string) (*Client, error) {
	apiClient, err := apiclient.NewClient(baseURL)
	if err != nil {
		return nil, err
	}
	client := &Client{apiClient: apiClient}
	return client, nil
}

func (client *Client) CheckName(ctx context.Context, name string) (bool, error) {
	return client.Check(ctx, CheckTypeName, name)
}

func (client *Client) CheckEmail(ctx context.Context, email string) (bool, error) {
	return client.Check(ctx, CheckTypeEmail, email)
}

func (client *Client) Check(ctx context.Context, checkType, content string) (bool, error) {
	data := url.Values{}
	data.Add(typeKey, checkType)
	data.Add(contentKey, content)

	resBody, err := client.doRequestWithRetries(ctx, checkEndpoint, data)
	if err != nil {
		return false, err
	}

	if string(resBody) == checkOk {
		return true, nil
	}

	return false, nil
}

func (client *Client) Match(ctx context.Context, checkType, content string) ([]MatchRecord, error) {
	response := make([]MatchRecord, 0)

	data := url.Values{}
	data.Add(typeKey, checkType)
	data.Add(contentKey, content)

	resBody, err := client.doRequestWithRetries(ctx, matchEndpoint, data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *Client) doRequestWithRetries(ctx context.Context, endpoint string, data url.Values) ([]byte, error) {
	var resBody []byte
	var err error
	tries := 1

	for tries <= retriesLimit {
		resBody, err = client.doRequest(ctx, endpoint, data)
		if err == nil {
			break
		}

		if tries < retriesLimit {
			time.Sleep(retrySleep * time.Millisecond)
		}
		tries++
	}

	return resBody, err
}

func (client *Client) doRequest(ctx context.Context, endpoint string, data url.Values) ([]byte, error) {
	headers := http.Header{}
	span, _ := opentracing.StartSpanFromContext(ctx, "phalanx-request")
	if span != nil {
		defer span.Finish()
		ext.SpanKindRPCClient.Set(span)
		ext.HTTPUrl.Set(span, endpoint)
		ext.HTTPMethod.Set(span, http.MethodPost)
		span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
	}

	resp, err := client.apiClient.Call("POST", endpoint, data, tracing.AddHttpHeadersFromContext(headers, ctx))
	if err != nil {
		return nil, err
	}
	
	if span != nil {
		ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))
	}

	resBody, err := client.apiClient.GetBody(resp)
	if err != nil {
		return nil, err
	}
	
	if span != nil {
		span.LogKV("response", string(resBody))
	}

	return resBody, nil
}
