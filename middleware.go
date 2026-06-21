package cafe

import "net/http"

type (
	Handler     = http.Handler
	HandlerFunc = http.HandlerFunc
)

type Middleware func(next HandlerFunc) HandlerFunc

func chain(f HandlerFunc, mws []Middleware) HandlerFunc {
	if len(mws) == 0 {
		return f
	}

	for i := len(mws) - 1; i >= 0; i-- {
		f = mws[i](f)
	}

	return f
}
