package quick

import (
	"fmt"
	"net/http"
	"strings"
)

// TODO: Cambiar los mapas por arrays -> orden determinista
type App struct {
	server  http.Server
	handler *http.ServeMux
	routers []mountedRouter
}

func NewServer() App {
	return App{
		handler: http.NewServeMux(),
		routers: []mountedRouter{},
	}
}

func (a *App) UseRouter(path string, ro *router) {
	for _, mr := range a.routers {
		if mr.path == path {
			return
		}
	}
	a.routers = append(a.routers, mountedRouter{path: path, router: ro})
}

func (a *App) Listen(addr string) error {
	a.setUpRouters()
	a.server = http.Server{
		Addr:    addr,
		Handler: a.handler,
	}
	return a.server.ListenAndServe()
}

func (a *App) setUpRouters() {
	for _, mr := range a.routers {
		routes := mr.router.getRoutes()
		for _, r := range routes {
			path := mr.path + r.path
			a.handle(path, r.method, r.handler)
		}
	}
}

func (a *App) handle(path, method string, handler http.HandlerFunc) {
	patt := fmt.Sprintf("%s %s", method, path)
	if !strings.HasSuffix(patt, "/") {
		patt += "/"
	}
	patt += "{$}"
	a.handler.HandleFunc(patt, handler)
}

/*** Basic HTTP Methods ***/

func (a *App) Get(path string, handler http.HandlerFunc) {
	a.handle(path, "GET", handler)
}

func (a *App) Post(path string, handler http.HandlerFunc) {
	a.handle(path, "POST", handler)
}

func (a *App) Put(path string, handler http.HandlerFunc) {
	a.handle(path, "PUT", handler)
}

func (a *App) Delete(path string, handler http.HandlerFunc) {
	a.handle(path, "DELETE", handler)
}
