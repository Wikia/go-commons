package consul

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	catalog := new(MockedConsulCatalogAPI)
	resolver := NewResolver(catalog)

	services, err := resolver.FindAll("auth", "production")
	assert.Equal(t, 2, len(services), "we did not receive 2 nodes")
	assert.Equal(t, "10.10.10.10", services[0].Address, "unexpected ip")
	assert.Equal(t, 9500, services[0].Port, "unexpected port")
	assert.Nil(t, err)
}

func TestResolve(t *testing.T) {
	catalog := new(MockedConsulCatalogAPI)
	resolver := NewResolver(catalog)

	service, err := resolver.Resolve("auth", "production")
	assert.Equal(t, "10.10.10.10", service.Address, "unexpected ip")
	assert.Equal(t, 9500, service.Port, "unexpected port")
	assert.Nil(t, err)
}
