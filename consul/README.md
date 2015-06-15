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
config you can create a `Health` client directly and inject that into `NewResolver`. Example:

```go
import (
	"github.com/hashicorp/consul/api"
	"github.com/Wikia/go-commons/consul"
)

config := api.DefaultConfig()
client, _ := api.NewClient(config)
health := client.Health()
resolver := consul.NewResolver(health)
address, _ := resolver.ResolveURI("user-preference", "production") // returns "http://10.10.10.10:12345"
```

If you only need to change the consul address, you can do set the following
environment variable and use `NewDefaultResolver()`:

```
export CONSUL_HTTP_ADDR=consul.service.consul:8500
```

See the package for more details regarding the configuration.
