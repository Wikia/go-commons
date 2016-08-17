package lars

import (
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// IRouteGroup interface for router group
type IRouteGroup interface {
	IRoutes
	Group(prefix string, middleware ...Handler) IRouteGroup
}

// IRoutes interface for routes
type IRoutes interface {
	Use(...Handler)
	Any(string, ...Handler)
	Get(string, ...Handler)
	Post(string, ...Handler)
	Delete(string, ...Handler)
	Patch(string, ...Handler)
	Put(string, ...Handler)
	Options(string, ...Handler)
	Head(string, ...Handler)
	Connect(string, ...Handler)
	Trace(string, ...Handler)
	WebSocket(websocket.Upgrader, string, Handler)
}

// routeGroup struct containing all fields and methods for use.
type routeGroup struct {
	prefix     string
	middleware HandlersChain
	lars       *LARS
}

var _ IRouteGroup = &routeGroup{}

func (g *routeGroup) handle(method string, path string, handlers []Handler) {

	if len(handlers) == 0 {
		panic("No handler mapped to path:" + path)
	}

	if i := strings.Index(path, "//"); i != -1 {
		panic("Bad path '" + path + "' contains duplicate // at index:" + strconv.Itoa(i))
	}

	chain := make(HandlersChain, len(handlers))
	name := ""

	for i, h := range handlers {

		if i == len(handlers)-1 {
			chain[i], name = g.lars.wrapHandlerWithName(h)
		} else {
			chain[i] = g.lars.wrapHandler(h)
		}
	}

	tree := g.lars.trees[method]
	if tree == nil {
		tree = new(node)
		g.lars.trees[method] = tree
	}

	combined := make(HandlersChain, len(g.middleware)+len(chain))
	copy(combined, g.middleware)
	copy(combined[len(g.middleware):], chain)

	pCount := tree.add(g.prefix+path, name, combined)
	pCount++

	if pCount > g.lars.mostParams {
		g.lars.mostParams = pCount
	}
}

// Use adds a middleware handler to the group middleware chain.
func (g *routeGroup) Use(m ...Handler) {
	for _, h := range m {
		g.middleware = append(g.middleware, g.lars.wrapHandler(h))
	}
}

// Connect adds a CONNECT route & handler to the router.
func (g *routeGroup) Connect(path string, h ...Handler) {
	g.handle(CONNECT, path, h)
}

// Delete adds a DELETE route & handler to the router.
func (g *routeGroup) Delete(path string, h ...Handler) {
	g.handle(DELETE, path, h)
}

// Get adds a GET route & handler to the router.
func (g *routeGroup) Get(path string, h ...Handler) {
	g.handle(GET, path, h)
}

// Head adds a HEAD route & handler to the router.
func (g *routeGroup) Head(path string, h ...Handler) {
	g.handle(HEAD, path, h)
}

// Options adds an OPTIONS route & handler to the router.
func (g *routeGroup) Options(path string, h ...Handler) {
	g.handle(OPTIONS, path, h)
}

// Patch adds a PATCH route & handler to the router.
func (g *routeGroup) Patch(path string, h ...Handler) {
	g.handle(PATCH, path, h)
}

// Post adds a POST route & handler to the router.
func (g *routeGroup) Post(path string, h ...Handler) {
	g.handle(POST, path, h)
}

// Put adds a PUT route & handler to the router.
func (g *routeGroup) Put(path string, h ...Handler) {
	g.handle(PUT, path, h)
}

// Trace adds a TRACE route & handler to the router.
func (g *routeGroup) Trace(path string, h ...Handler) {
	g.handle(TRACE, path, h)
}

// Handle allows for any method to be registered with the given
// route & handler. Allows for non standard methods to be used
// like CalDavs PROPFIND and so forth.
func (g *routeGroup) Handle(method string, path string, h ...Handler) {
	g.handle(method, path, h)
}

// Any adds a route & handler to the router for all HTTP methods.
func (g *routeGroup) Any(path string, h ...Handler) {
	g.Connect(path, h...)
	g.Delete(path, h...)
	g.Get(path, h...)
	g.Head(path, h...)
	g.Options(path, h...)
	g.Patch(path, h...)
	g.Post(path, h...)
	g.Put(path, h...)
	g.Trace(path, h...)
}

// Match adds a route & handler to the router for multiple HTTP methods provided.
func (g *routeGroup) Match(methods []string, path string, h ...Handler) {
	for _, m := range methods {
		g.handle(m, path, h)
	}
}

// WebSocket adds a websocket route
func (g *routeGroup) WebSocket(upgrader websocket.Upgrader, path string, h Handler) {

	handler := g.lars.wrapHandler(h)
	g.Get(path, func(c Context) {

		ctx := c.BaseContext()
		var err error

		ctx.websocket, err = upgrader.Upgrade(ctx.response, ctx.request, nil)
		if err != nil {
			return
		}

		defer ctx.websocket.Close()
		c.Next()
	}, handler)
}

// Group creates a new sub router with prefix. It inherits all properties from
// the parent. Passing middleware overrides parent middleware but still keeps
// the root level middleware intact.
func (g *routeGroup) Group(prefix string, middleware ...Handler) IRouteGroup {

	rg := &routeGroup{
		prefix: g.prefix + prefix,
		lars:   g.lars,
	}

	if len(middleware) == 0 {
		rg.middleware = make(HandlersChain, len(g.middleware)+len(middleware))
		copy(rg.middleware, g.middleware)

		return rg
	}

	if middleware[0] == nil {
		rg.middleware = make(HandlersChain, 0)
		return rg
	}

	rg.middleware = make(HandlersChain, len(middleware))
	copy(rg.middleware, g.lars.middleware)
	rg.Use(middleware...)

	return rg
}
