package phalanx

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/mock"
)

type ApiClientMock struct {
	mock.Mock
}

func (m *ApiClientMock) SetRetryMax(count int) {
	m.Called(count)
}

func (m *ApiClientMock) SetRetryWaitMin(period time.Duration) {
	m.Called(period)
}

func (m *ApiClientMock) SetRetryWaitMax(period time.Duration) {
	m.Called(period)
}

func (m *ApiClientMock) SetLogger(log *log.Logger) {
	m.Called(log)
}

func (m *ApiClientMock) GetClient() *retryablehttp.Client {
	args := m.Called()

	return args.Get(0).(*retryablehttp.Client)
}

func (m *ApiClientMock) CallWithContext(ctx context.Context, method, endpoint string, data url.Values, headers http.Header) (*http.Response, error) {
	args := m.Called(ctx, method, endpoint, data, headers)

	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *ApiClientMock) Call(method, endpoint string, data url.Values, headers http.Header) (*http.Response, error) {
	args := m.Called(method, endpoint, data, headers)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *ApiClientMock) GetBody(resp *http.Response) ([]byte, error) {
	args := m.Called(resp)
	return args.Get(0).([]byte), args.Error(1)
}
