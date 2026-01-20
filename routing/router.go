package quick

import (
	"fmt"
	"net/http"
	"strings"
)

type route struct {
	url     string
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
		routes:  make([]route, 0),
		routers: make(map[string]*router, 0),
	}
}

func (r *router) UseRouter(path string, ro *router) {
	r.routers[path] = ro
}

func (r *router) getRoutes() map[string]http.HandlerFunc {
	mountedRoutes := make(map[string]http.HandlerFunc)
	for _, ro := range r.routes {
		patt := fmt.Sprintf("%s %s", ro.method, ro.url)
		mountedRoutes[patt] = ro.handler
	}
	for path, rtr := range r.routers {
		rtrRoutes := rtr.getRoutes()
		for k, v := range rtrRoutes {
			parts := strings.Split(k, " ")
			patt := fmt.Sprintf("%s %s%s", parts[0], path, parts[1])
			mountedRoutes[patt] = v
		}
	}
	return mountedRoutes
}

/*** Basic HTTP Methods ***/

func (r *router) Get(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "GET",
		url:     url,
		handler: handler,
	})
}

func (r *router) Post(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "POST",
		url:     url,
		handler: handler,
	})
}

func (r *router) Put(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "PUT",
		url:     url,
		handler: handler,
	})
}

func (r *router) Delete(url string, handler http.HandlerFunc) {
	r.routes = append(r.routes, route{
		method:  "DELETE",
		url:     url,
		handler: handler,
	})
}
