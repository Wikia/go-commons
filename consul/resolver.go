package consul

import (
	"fmt"
	"math/rand"

	"github.com/hashicorp/consul/api"
)

type ConsulCatalogAPI interface {
	Service(service, tag string, passingOnly bool, q *api.QueryOptions) ([]*api.ServiceEntry,
		*api.QueryMeta, error)
}

type AddressTuple struct {
	Address string
	Port    int
}

type ConsulResolver interface {
	ResolveAll(name, tag string) ([]*AddressTuple, error)
	Resolve(name, tag string) (*AddressTuple, error)
	ResolveAddress(name, tag string) (string, error)
}

type ConsulResolverValue struct {
	consulAPIClient ConsulCatalogAPI
}

func NewResolver(consulAPIClient ConsulCatalogAPI) *ConsulResolverValue {
	return &ConsulResolverValue{consulAPIClient}
}

func (resolver *ConsulResolverValue) ResolveAll(name, tag string) ([]*AddressTuple, error) {
	services, _, err := resolver.consulAPIClient.Service(name, tag, true, nil)
	if err != nil {
		return nil, err
	}

	tuples := make([]*AddressTuple, len(services))
	for index, service := range services {
		tuples[index] = &AddressTuple{service.Node.Address, service.Service.Port}
	}

	return tuples, nil
}

func (resolver *ConsulResolverValue) Resolve(name, tag string) (*AddressTuple, error) {
	services, err := resolver.ResolveAll(name, tag)
	if err != nil {
		return nil, err
	}

	if services != nil && len(services) < 1 {
		return nil, fmt.Errorf("error, no services available for %q with tag %q", name, tag)
	}

	return services[rand.Intn(len(services))], nil
}

func (resolver *ConsulResolverValue) ResolveAddress(name, tag string) (string, error) {
	return "", nil
}

func (t *AddressTuple) ToAddress() string {
	if t != nil {
		return fmt.Sprintf("http://%s:%d", t.Address, t.Port)
	}

	return ""
}
