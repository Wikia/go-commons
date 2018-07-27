package apiclient

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type ApiClient interface {
	SetLogger(log *log.Logger)
	SetRetryMax(count int)
	SetRetryWaitMin(period time.Duration)
	SetRetryWaitMax(period time.Duration)
	CallWithContext(ctx context.Context, method, endpoint string, data url.Values, headers http.Header) (*http.Response, error)
	Call(method, endpoint string, data url.Values, headers http.Header) (*http.Response, error)
	NewRequest(method, endpoint string, data url.Values) (*retryablehttp.Request, error)
	GetBody(resp *http.Response) ([]byte, error)
	GetClient() *retryablehttp.Client
}

type Client struct {
	httpClient *retryablehttp.Client
	BaseURL    *url.URL
}

func NewClient(baseURL string) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	client := &Client{httpClient: retryablehttp.NewClient(), BaseURL: parsedURL}
	client.SetRetryMax(3)
	client.SetRetryWaitMin(20 * time.Nanosecond)
	client.SetRetryWaitMin(1 * time.Second)

	return client, nil
}

func NewClientWithProxy(baseURL string, proxy string) (*Client, error) {
	client, err := NewClient(baseURL)
	if proxy == "" || err != nil {
		return client, err
	}

	proxyURL, err := url.Parse(proxy)
	if err == nil {
		client.httpClient.HTTPClient.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	return client, nil
}

func (client *Client) GetClient() *retryablehttp.Client {
	return client.httpClient
}

func (client *Client) SetLogger(log *log.Logger) {
	client.httpClient.Logger = log
}

func (client *Client) SetRetryMax(count int) {
	client.httpClient.RetryMax = count
}

func (client *Client) SetRetryWaitMin(period time.Duration) {
	client.httpClient.RetryWaitMin = period
}

func (client *Client) SetRetryWaitMax(period time.Duration) {
	client.httpClient.RetryWaitMax = period
}

func (client *Client) CallWithContext(ctx context.Context, method, endpoint string, data url.Values, headers http.Header) (*http.Response, error) {
	request, err := client.NewRequest(method, endpoint, data)

	span, newCtx := opentracing.StartSpanFromContext(ctx, "apiClient_call")
	if span != nil {
		defer span.Finish()

		ext.SpanKindRPCClient.Set(span)
		ext.HTTPMethod.Set(span, method)
		ext.HTTPUrl.Set(span, endpoint)
		span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
	}

	if err != nil {
		if span != nil {
			ext.Error.Set(span, true)
			span.LogKV("reason", err)
		}
		return nil, err
	}

	if newCtx != nil {
		request = request.WithContext(newCtx)
	}

	// adding headers
	for key, values := range headers {
		for _, v := range values {
			request.Header.Set(key, v)
		}
	}

	// This seems heavy handed but as a rule we are closing the connection after
	// GetBody below. This ensures that we are communicating our intentions in
	// the HTTP request.
	request.Close = true

	response, err := client.httpClient.Do(request)
	if err != nil {
		if span != nil {
			ext.Error.Set(span, true)
			span.LogKV("reason", err)
		}
		return nil, err
	}

	if span != nil {
		ext.HTTPStatusCode.Set(span, uint16(response.StatusCode))
	}

	return response, nil
}

func (client *Client) Call(method, endpoint string, data url.Values, headers http.Header) (*http.Response, error) {
	request, err := client.NewRequest(method, endpoint, data)
	if err != nil {
		return nil, err
	}

	// adding headers
	for key, values := range headers {
		for _, v := range values {
			request.Header.Set(key, v)
		}
	}

	// This seems heavy handed but as a rule we are closing the connection after
	// GetBody below. This ensures that we are communicating our intentions in
	// the HTTP request.
	request.Close = true

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (client *Client) NewRequest(method, endpoint string, data url.Values) (*retryablehttp.Request, error) {
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	fullUrl := client.BaseURL.ResolveReference(endpointUrl)

	request, err := retryablehttp.NewRequest(method, fullUrl.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	return request, nil
}

func (client *Client) GetBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
