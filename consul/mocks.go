package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/mock"
)

type MockedConsulHealthAPI struct {
	mock.Mock
}

func (m *MockedConsulHealthAPI) Service(service, tag string, passingOnly bool, q *api.QueryOptions) ([]*api.ServiceEntry, *api.QueryMeta, error) {
	args := m.Called(service, tag, passingOnly, q)
	return args.Get(0).([]*api.ServiceEntry), args.Get(1).(*api.QueryMeta), args.Error(2)
}
