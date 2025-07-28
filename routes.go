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
