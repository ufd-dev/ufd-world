package main

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var tpl *template.Template

func configRoutes() *mux.Router {
	tpl = template.Must(template.ParseGlob("templates/*.tpl.html"))

	r := mux.NewRouter()
	htmlRouter := r.NewRoute().Subrouter()
	htmlRouter.Use(contentTypeHTMLMiddleware)

	r.HandleFunc("/", handleTempWelcome)
	return r
}

func contentTypeHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func handleTempWelcome(w http.ResponseWriter, req *http.Request) {
	err := tpl.ExecuteTemplate(w, "main.tpl.html", nil)
	if err != nil {
		w.Write([]byte("An unknown error has occurred."))
	}
}
