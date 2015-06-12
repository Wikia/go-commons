# Consul Package

This package provides a high level interface for resolving the URI for a
service that's registered in consul.

## Usage

Here is a basic example that should suffice for most cases:

```go

import (
	"github.com/Wikia/go-commons/consul"
)

resolver := consul.NewDefaultResolver()
address, _ := resolver.ResolveURI("user-preference", "production") // returns "http://10.10.10.10:12345"
```

The above uses some sane defaults for consul. If you need to use a different
host other than `consul.service.consul` you can create a `Health` client
directly and inject that into `NewResolver`. Example:

```go
import (
	"github.com/Wikia/go-commons/consul"
)

health := consul.ConsulAPIHealthClientFactory(consul.ConsulAPIConfigFactory("your.consul.service:8500"))
resolver := consul.NewResolver(health)
address, _ := resolver.ResolveURI("user-preference", "production") // returns "http://10.10.10.10:12345"
```

See the package for more details regarding the configuration.
