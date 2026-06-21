package middleware

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/LucasSim0n/cafe"
	"github.com/LucasSim0n/cafe/internal/httpx"
)

type LoggerConfig struct {
	Output io.Writer
}

func Logger(cfg LoggerConfig) cafe.Middleware {
	return func(next cafe.HandlerFunc) cafe.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			rw := httpx.NewResponseWriter(w)

			start := time.Now()

			next(rw, r)

			fmt.Fprintf(cfg.Output,
				"%s %s %d %s",
				r.Method,
				r.URL.Path,
				rw.Status,
				time.Since(start),
			)

		}
	}
}
