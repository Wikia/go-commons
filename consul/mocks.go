package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/mock"
)

type MockedConsulCatalogAPI struct {
	mock.Mock
}

func (m *MockedConsulCatalogAPI) Service(service, tag string, q *api.QueryOptions) ([]*api.CatalogService, *api.QueryMeta, error) {
	args := m.Called(service, tag, q)
	return args.Get(0).([]*api.CatalogService), args.Get(1).(*api.QueryMeta), args.Error(2)
}
