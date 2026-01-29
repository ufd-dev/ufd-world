package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ufd-dev/ufd-world/media"
)

func configRoutes() *mux.Router {
	r := mux.NewRouter()

	// static files
	fs := http.FileServer(http.Dir("./static/"))
	r.PathPrefix("/static").Handler(
		http.StripPrefix(
			"/static/",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// max-age=1year
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
				fs.ServeHTTP(w, r)
			})))

	// HTML pages
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		renderTemplate(w, "home.tpl.html", nil)
	})
	r.HandleFunc("/media", func(w http.ResponseWriter, req *http.Request) {
		renderTemplate(w, "media.tpl.html", nil)
	})
	r.HandleFunc("/img-tagger", func(w http.ResponseWriter, req *http.Request) {
		renderTemplate(w, "img-tagger.tpl.html", nil)
	})

	s := r.PathPrefix("/img-tagger/download").Subrouter()
	s.Use(noCacheMiddleware)
	s.HandleFunc("", handleDownloadTaggedImg)

	// JSON API
	ar := r.PathPrefix("/api").Subrouter()
	ar.Use(contentTypeJSONMiddleware)
	ar.HandleFunc("/media", func(w http.ResponseWriter, req *http.Request) {
		media, err := media.GetList()
		if err != nil {
			http.Error(w, "\"Internal server error\"", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
		err = json.NewEncoder(w).Encode(media)
		if err != nil {
			http.Error(w, "\"Internal server error\"", http.StatusInternalServerError)
			fmt.Println(err)
			return
		}
	})

	return r
}

func contentTypeJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}
