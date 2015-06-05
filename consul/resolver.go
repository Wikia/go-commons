package consul

import (
	"fmt"
	"math/rand"

	"github.com/hashicorp/consul/api"
)

type ConsulCatalogAPI interface {
	Service(service, tag string, q *api.QueryOptions) ([]*api.CatalogService, *api.QueryMeta, error)
}

type AddressTuple struct {
	Address string
	Port    int
}

type ConsulResolver interface {
	FindAll(name, tag string) (*[]AddressTuple, error)
	Resolve(name, tag string) (*AddressTuple, error)
}

type ConsulResolverValue struct {
	consulAPIClient ConsulCatalogAPI
}

func NewResolver(consulAPIClient ConsulCatalogAPI) *ConsulResolverValue {
	return &ConsulResolverValue{consulAPIClient}
}

func (resolver *ConsulResolverValue) FindAll(name, tag string) ([]*AddressTuple, error) {
	return make([]*AddressTuple, 0), nil
}

func (resolver *ConsulResolverValue) Resolve(name, tag string) (*AddressTuple, error) {
	services, err := resolver.FindAll(name, tag)
	if err != nil {
		return nil, err
	}

	if services != nil && len(services) < 1 {
		return nil, fmt.Errorf("error, empty service list for %q with tag %q", name, tag)
	}

	return services[rand.Intn(len(services))], nil
}
