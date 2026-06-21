package cafe

import (
	"fmt"
	"net/http"
	"strings"
)

type App struct {
	server      http.Server
	handler     *http.ServeMux
	routers     []mountedRouter
	routes      []route
	middlewares []Middleware
}

/*** Factory ***/

func NewServer() App {
	return App{
		handler:     http.NewServeMux(),
		routers:     []mountedRouter{},
		routes:      []route{},
		middlewares: []Middleware{},
	}
}

/*** Aggregation ***/

func (a *App) UseRouter(path string, ro *Router) {
	for _, mr := range a.routers {
		if mr.path == path {
			return
		}
	}
	a.routers = append(a.routers, mountedRouter{path: path, router: ro})
}

func (a *App) Use(mw Middleware) {
	a.middlewares = append(a.middlewares, mw)
}

/*** Setup ***/

func (a *App) Listen(addr string) error {
	a.SetUpRouters()
	a.server = http.Server{
		Addr:    addr,
		Handler: a.handler,
	}
	return a.server.ListenAndServe()
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}

/*** Basic HTTP Methods ***/

func (a *App) Get(path string, handler HandlerFunc) {
	a.routes = addRoute(a.routes, path, "GET", handler)
}

func (a *App) Post(path string, handler HandlerFunc) {
	a.routes = addRoute(a.routes, path, "POST", handler)
}

func (a *App) Put(path string, handler HandlerFunc) {
	a.routes = addRoute(a.routes, path, "PUT", handler)
}

func (a *App) Delete(path string, handler HandlerFunc) {
	a.routes = addRoute(a.routes, path, "DELETE", handler)
}

/*** Utils ***/

func (a *App) SetUpRouters() {
	for _, r := range a.routes {
		h := chain(r.handler, a.middlewares)
		a.handle(r.path, r.method, h)
	}

	for _, mr := range a.routers {
		routes := mr.router.getRoutes()
		for _, r := range routes {
			path := mr.path + r.path
			h := chain(r.handler, a.middlewares)
			a.handle(path, r.method, h)
		}
	}
}

func (a *App) handle(path, method string, handler HandlerFunc) {
	patt := fmt.Sprintf("%s %s", method, path)
	if !strings.HasSuffix(patt, "/") {
		patt += "/"
	}
	patt += "{$}"
	a.handler.Handle(patt, handler)
}

func addRoute(routes []route, path, method string, handler HandlerFunc) []route {
	for _, r := range routes {
		if r.path == path && r.method == method {
			return routes
		}
	}

	return append(routes, route{
		path:    path,
		method:  method,
		handler: handler,
	})
}
