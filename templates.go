package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

var templates map[string]*template.Template

func loadTemplates() {
	templates = make(map[string]*template.Template)

	mainTpl := template.Must(template.ParseGlob("templates/layouts/*.tpl.html"))

	pagePaths, err := filepath.Glob("./templates/*.tpl.html")
	if err != nil {
		panic("cannot load page templates")
	}

	for _, pp := range pagePaths {
		file := filepath.Base(pp)
		templates[file] = template.Must(template.Must(mainTpl.Clone()).ParseFiles("templates/" + file))
	}
}

func renderTemplate(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, ok := templates[name]
	if !ok {
		errTxt := fmt.Sprintf("404: %s not found", name)
		http.Error(w, errTxt, http.StatusInternalServerError)
		fmt.Println(errTxt)
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println("error rendering", name, err)
	}
}
