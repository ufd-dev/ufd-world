package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func configRoutes() *mux.Router {
	r := mux.NewRouter()

	// HTML pages
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		renderTemplate(w, "home.tpl.html", nil)
	})
	r.HandleFunc("/media", func(w http.ResponseWriter, req *http.Request) {
		renderTemplate(w, "media.tpl.html", nil)
	})

	// JSON API
	ar := r.PathPrefix("/api").Subrouter()
	ar.Use(contentTypeJSONMiddleware)
	ar.HandleFunc("/media", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("[]"))
	})

	return r
}

func contentTypeJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
