package main

import (
	"log"
	"net/http"

	quick "github.com/LucasSim0n/quick/routing"
)

func main() {
	app := quick.NewServer()
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello World!"))
	})

	r := quick.NewRouter()
	r.Get("/myproblem", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello from router"))
	})

	r2 := quick.NewRouter()
	r2.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello from embeded router"))
	})

	r.UseRouter("/r2", r2)

	app.UseRouter("/router", r)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatal(err)
	}

}
