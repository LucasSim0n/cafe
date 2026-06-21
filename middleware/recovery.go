package middleware

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"

	"github.com/LucasSim0n/cafe"
	"github.com/LucasSim0n/cafe/internal/httpx"
)

type RecoveryConfig struct {
	PrintStack bool
	Output     io.Writer
}

func Recovery(cfg RecoveryConfig) cafe.Middleware {
	return func(next cafe.HandlerFunc) cafe.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			rw := httpx.NewResponseWriter(w)

			defer func() {
				if err := recover(); err != nil {
					if cfg.PrintStack {
						fmt.Fprintf(cfg.Output, "panic recovered: %v\n%s", err, debug.Stack())
					} else {
						fmt.Fprintf(cfg.Output, "panic recovered: %v\n", err)
					}
				}

				if !rw.ResponseStarted {
					http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next(rw, r)
		}
	}
}
