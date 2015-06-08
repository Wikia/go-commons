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
	ResolveAll(name, tag string, allowStale bool) ([]*AddressTuple, error)
	Resolve(name, tag string, allowStale bool) (*AddressTuple, error)
	ResolveAddress(name, tag string, allowStale bool) (string, error)
}

type ConsulResolverValue struct {
	consulAPIClient ConsulCatalogAPI
}

func NewResolver(consulAPIClient ConsulCatalogAPI) *ConsulResolverValue {
	return &ConsulResolverValue{consulAPIClient}
}

func DefaultResolver() *ConsulResolverValue {
	config := api.DefaultConfig()
	config.Address = "consul.service.consul:8500"
	client, _ := api.NewClient(config)
	health := client.Health()
	return &ConsulResolverValue{health}
}

func (resolver *ConsulResolverValue) ResolveAll(name, tag string, allowStale bool) ([]*AddressTuple, error) {
	options := &api.QueryOptions{AllowStale: allowStale}
	services, _, err := resolver.consulAPIClient.Service(name, tag, true, options)
	if err != nil {
		return nil, err
	}

	tuples := make([]*AddressTuple, len(services))
	for index, service := range services {
		tuples[index] = &AddressTuple{service.Node.Address, service.Service.Port}
	}

	return tuples, nil
}

func (resolver *ConsulResolverValue) Resolve(name, tag string, allowStale bool) (*AddressTuple, error) {
	services, err := resolver.ResolveAll(name, tag, allowStale)
	if err != nil {
		return nil, err
	}

	if services != nil && len(services) < 1 {
		return nil, fmt.Errorf("error, no services available for %q with tag %q", name, tag)
	}

	return services[rand.Intn(len(services))], nil
}

func (resolver *ConsulResolverValue) ResolveAddress(name, tag string, allowStale bool) (string, error) {
	tuple, err := resolver.Resolve(name, tag, allowStale)
	if err != nil {
		return "", err
	}

	return tuple.ToAddress(), nil
}

func (t *AddressTuple) ToAddress() string {
	if t != nil {
		return fmt.Sprintf("http://%s:%d", t.Address, t.Port)
	}

	return ""
}
