package quick

import (
	"net/http"
)

type route struct {
	path    string
	method  string
	handler http.HandlerFunc
}

type router struct {
	url     string
	routes  []route
	routers map[string]*router
}

func NewRouter() *router {
	return &router{
		routes:  []route{},
		routers: map[string]*router{},
	}
}

func (r *router) UseRouter(path string, ro *router) {
	r.routers[path] = ro
}

func (r *router) getRoutes() []route {
	mountedRoutes := r.routes[:]
	for path, rtr := range r.routers {
		rtrRoutes := rtr.getRoutes()
		for _, v := range rtrRoutes {
			v.path = path + v.path
			mountedRoutes = append(mountedRoutes, v)
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

func (r *router) Post(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "POST",
		path:    url,
		handler: handler,
	})
}

func (r *router) Put(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "PUT",
		path:    url,
		handler: handler,
	})
}

func (r *router) Delete(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "DELETE",
		path:    url,
		handler: handler,
	})
}
