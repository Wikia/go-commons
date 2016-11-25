package phalanx

import (
	"net/http"
	"net/url"

	"github.com/stretchr/testify/mock"
)

type ApiClientMock struct {
	mock.Mock
}

func (m *ApiClientMock) Call(method, endpoint string, data url.Values, headers map[string]string) (*http.Response, error) {
	args := m.Called(method, endpoint, data, headers)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *ApiClientMock) NewRequest(method, endpoint string, data url.Values) (*http.Request, error) {
	args := m.Called(method, endpoint, data)
	return args.Get(0).(*http.Request), args.Error(1)
}

func (m *ApiClientMock) GetBody(resp *http.Response) ([]byte, error) {
	args := m.Called(resp)
	return args.Get(0).([]byte), args.Error(1)
}
