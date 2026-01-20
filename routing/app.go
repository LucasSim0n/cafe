package quick

import (
	"fmt"
	"net/http"
	"strings"
)

type App struct {
	server  http.Server
	handler *http.ServeMux
	routers map[string]*router
}

func NewServer() App {
	return App{
		handler: http.NewServeMux(),
		routers: make(map[string]*router, 0),
	}
}

func (a *App) UseRouter(url string, r *router) {
	a.routers[url] = r
}

func (a *App) Listen(addr string) {
	a.setUpRouters()
	a.server = http.Server{
		Addr:    addr,
		Handler: a.handler,
	}
	a.server.ListenAndServe()
}

func (a *App) setUpRouters() {
	for p, ro := range a.routers {
		rts := ro.getRoutes()
		for k, v := range rts {
			parts := strings.Split(k, " ")
			patt := fmt.Sprintf("%s %s%s", parts[0], p, parts[1])
			a.handler.HandleFunc(patt, v)
			fmt.Println(patt)
		}
	}
}

/*** Basic HTTP Methods ***/

func (a *App) Get(url string, handler http.HandlerFunc) {
	patt := fmt.Sprintf("GET %s", url)
	a.handler.HandleFunc(patt, handler)
}

func (a *App) Post(url string, handler http.HandlerFunc) {
	patt := fmt.Sprintf("POST %s", url)
	a.handler.HandleFunc(patt, handler)
}

func (a *App) Put(url string, handler http.HandlerFunc) {
	patt := fmt.Sprintf("PUT %s", url)
	a.handler.HandleFunc(patt, handler)
}

func (a *App) Delete(url string, handler http.HandlerFunc) {
	patt := fmt.Sprintf("DELETE %s", url)
	a.handler.HandleFunc(patt, handler)
}
