package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/LucasSim0n/cafe"
)

func TestPanicRecoveryMiddleware(t *testing.T) {
	app := cafe.NewServer()
	var info strings.Builder
	app.Use(Recovery(
		RecoveryConfig{
			Output: &info,
		},
	))

	app.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
		panic("Testing")
	})

	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	httptest.NewRequest("GET", "/panic", nil)
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK after panic recovery, got %d", rr.Code)
	}
}

func TestLoggerMiddleware(t *testing.T) {
	app := cafe.NewServer()
	var info strings.Builder

	app.Use(Logger(LoggerConfig{Output: &info}))
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK after panic recovery, got %d", rr.Code)
	}

	if !strings.HasPrefix(info.String(), "GET") {
		t.Errorf("Expected method in log, but got %s", info.String())
	}
}
