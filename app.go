package cafe

import (
	"fmt"
	"net/http"
	"strings"
)

type App struct {
	*Router

	server http.Server
}

/*** Factory ***/

func NewServer() App {
	return App{
		Router: NewRouter(""),
	}
}

/*** Setup ***/

func (a *App) Listen(addr string) error {
	a.server = http.Server{
		Addr:    addr,
		Handler: a.buildMux(),
	}
	return a.server.ListenAndServe()
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.buildMux().ServeHTTP(w, r)
}

func (a *App) buildMux() *http.ServeMux {
	mux := http.NewServeMux()
	for _, r := range a.compileRoutes() {
		mux.Handle(buildHandlerString(r.method, r.path), r.handler)
	}

	return mux
}

/*** Utils ***/

func buildHandlerString(method, path string) string {
	patt := fmt.Sprintf("%s %s", method, path)
	if !strings.HasSuffix(patt, "/") {
		patt += "/"
	}
	patt += "{$}"
	return patt
}
