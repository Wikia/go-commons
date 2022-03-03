package consul

import (
	"math/rand"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ResolverTestSuite struct {
	suite.Suite
	consul   *MockedConsulHealthAPI
	resolver ConsulResolver
	response []*api.ServiceEntry
}

func (suite *ResolverTestSuite) SetupTest() {
	rand.Seed(2)
	suite.consul = new(MockedConsulHealthAPI)
	suite.resolver = NewResolver(suite.consul)

	suite.response = make([]*api.ServiceEntry, 2)
	suite.response[0] = &api.ServiceEntry{
		Node:    &api.Node{Node: "foo", Address: "10.10.10.10"},
		Service: &api.AgentService{ID: "1234", Service: "auth", Tags: []string{"production"}, Port: 9500, Address: ""},
		Checks:  nil,
	}
	suite.response[1] = &api.ServiceEntry{
		Node:    &api.Node{Node: "bar", Address: "10.10.10.11"},
		Service: &api.AgentService{ID: "1234", Service: "auth", Tags: []string{"production"}, Port: 9500, Address: ""},
		Checks:  nil,
	}
}

func (suite *ResolverTestSuite) TestResolveAll() {
	suite.consul.On("Service", "auth", "production", true,
		&api.QueryOptions{AllowStale: true}).Return(suite.response,
		(*api.QueryMeta)(nil), nil)
	services, err := suite.resolver.ResolveAll("auth", "production")
	assert.Equal(suite.T(), 2, len(services), "we did not receive 2 nodes")
	assert.Equal(suite.T(), "10.10.10.10", services[0].Address, "unexpected ip")
	assert.Equal(suite.T(), 9500, services[0].Port, "unexpected port")
	assert.Nil(suite.T(), err)
}

func (suite *ResolverTestSuite) TestResolve() {
	suite.consul.On("Service", "auth", "production", true,
		&api.QueryOptions{AllowStale: true}).Return(suite.response,
		(*api.QueryMeta)(nil), nil)
	service, err := suite.resolver.Resolve("auth", "production")
	assert.Equal(suite.T(), "10.10.10.10", service.Address, "unexpected ip")
	assert.Equal(suite.T(), 9500, service.Port, "unexpected port")
	assert.Nil(suite.T(), err)
}

func TestResolverTestSuite(t *testing.T) {
	suite.Run(t, new(ResolverTestSuite))
}

func TestAddressTupleToURI(t *testing.T) {
	tuple := &AddressTuple{"10.10.10.10", 80}
	address := tuple.ToURI()
	assert.Equal(t, "http://10.10.10.10:80", address, "error, malformed address")
}
