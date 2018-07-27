package phalanx

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/mock"
)

type ApiClientMock struct {
	mock.Mock
}

func (m *ApiClientMock) SetLogger(log *log.Logger) {
	m.Called(log)
}

func (m *ApiClientMock) CallWithContext(ctx context.Context, method, endpoint string, data url.Values, headers http.Header) (*http.Response, error) {
	args := m.Called(ctx, method, endpoint, data, headers)

	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *ApiClientMock) Call(method, endpoint string, data url.Values, headers http.Header) (*http.Response, error) {
	args := m.Called(method, endpoint, data, headers)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *ApiClientMock) NewRequest(method, endpoint string, data url.Values) (*retryablehttp.Request, error) {
	args := m.Called(method, endpoint, data)
	return args.Get(0).(*retryablehttp.Request), args.Error(1)
}

func (m *ApiClientMock) GetBody(resp *http.Response) ([]byte, error) {
	args := m.Called(resp)
	return args.Get(0).([]byte), args.Error(1)
}
