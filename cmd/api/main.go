package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type product struct {
	ID   int
	Name string
}

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		println("endpoint")

	})
	r.Get("/{productName}/{productId}", func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "productName")
		w.Write([]byte(param))
	})
	r.Get("/json", func(w http.ResponseWriter, r *http.Request) {
		obj := map[string]string{"message": "sucess"}
		render.JSON(w, r, obj)
	})
	r.Put("/product", func(w http.ResponseWriter, r *http.Request) {
		var product product
		render.DecodeJSON(r.Body, &product)
		product.ID = 5
		render.JSON(w, r, product)
	})

	http.ListenAndServe(":3000", r)
}

func myMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("before")
		next.ServeHTTP(w, r)
		println("after")
	})
}

func myMiddleware2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		println("request ", r.Method, " url ", r.RequestURI)
		next.ServeHTTP(w, r)

	})
}
