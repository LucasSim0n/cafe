package cafe

/*** Definitions ***/

type Router struct {
	prefix      string
	parent      *Router
	routes      []*Route
	children    []*Router
	middlewares []Middleware
}

type Route struct {
	path    string
	method  string
	handler HandlerFunc
}

/*** Init ***/

func NewRouter(prefix string) *Router {
	return &Router{
		prefix:      prefix,
		routes:      []*Route{},
		children:    []*Router{},
		middlewares: []Middleware{},
	}
}

/*** Aggregation ***/

func (r *Router) Group(path string) *Router {
	for _, c := range r.children {
		if c.prefix == path {
			return c
		}
	}

	child := &Router{
		prefix: path,
	}
	r.children = append(r.children, child)
	return child
}

func (r *Router) UseRouter(child *Router) {
	for _, c := range r.children {
		if c.prefix == child.prefix {
			return
		}
	}

	r.children = append(r.children, child)
}

func (r *Router) Use(mw Middleware) {
	r.middlewares = append(r.middlewares, mw)
}

/*** Basic HTTP Methods ***/

func (r *Router) Get(path string, handler HandlerFunc) {
	r.routes = addRoute(r.routes, path, "GET", handler)
}

func (r *Router) Post(path string, handler HandlerFunc) {
	r.routes = addRoute(r.routes, path, "POST", handler)
}

func (r *Router) Put(path string, handler HandlerFunc) {
	r.routes = addRoute(r.routes, path, "PUT", handler)
}

func (r *Router) Delete(path string, handler HandlerFunc) {
	r.routes = addRoute(r.routes, path, "DELETE", handler)
}

/*** Utils ***/

func (r *Router) compileRoutes() []*Route {
	compiledRoutes := []*Route{}

	for _, rt := range r.routes {
		comp := &Route{
			handler: chain(rt.handler, r.middlewares),
			path:    rt.path,
			method:  rt.method,
		}
		compiledRoutes = append(compiledRoutes, comp)
	}

	for _, child := range r.children {
		childRoutes := child.compileRoutes()
		for _, rt := range childRoutes {
			comp := &Route{
				path:    child.prefix + rt.path,
				method:  rt.method,
				handler: chain(rt.handler, r.middlewares),
			}
			compiledRoutes = append(compiledRoutes, comp)
		}
	}

	return compiledRoutes
}

func addRoute(routes []*Route, path, method string, handler HandlerFunc) []*Route {
	for _, r := range routes {
		if r.path == path && r.method == method {
			return routes
		}
	}

	return append(routes, &Route{
		path:    path,
		method:  method,
		handler: handler,
	})
}
