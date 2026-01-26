package quick

import (
	"net/http"
)

type route struct {
	path    string
	method  string
	handler http.HandlerFunc
}

// TODO: Cambiar los mapas por arrays -> orden determinista
type router struct {
	routes  []route
	routers []mountedRouter
}

type mountedRouter struct {
	path   string
	router *router
}

func NewRouter() *router {
	return &router{
		routes:  []route{},
		routers: []mountedRouter{},
	}
}

func (r *router) UseRouter(path string, ro *router) {
	for _, mr := range r.routers {
		if mr.path == path {
			return
		}
	}
	r.routers = append(r.routers, mountedRouter{path: path, router: ro})
}

func (r *router) getRoutes() []route {
	mountedRoutes := append([]route{}, r.routes...)
	for _, mr := range r.routers {
		rtrRoutes := mr.router.getRoutes()
		for _, rt := range rtrRoutes {
			rt.path = mr.path + rt.path
			mountedRoutes = append(mountedRoutes, rt)
		}
	}
	return mountedRoutes
}

/*** Basic HTTP Methods ***/

func (r *router) Get(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "GET",
		path:    path,
		handler: handler,
	})
}

func (r *router) Post(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "POST",
		path:    path,
		handler: handler,
	})
}

func (r *router) Put(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "PUT",
		path:    path,
		handler: handler,
	})
}

func (r *router) Delete(path string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "DELETE",
		path:    path,
		handler: handler,
	})
}
