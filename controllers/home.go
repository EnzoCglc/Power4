package controllers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render(w, "index.html", nil)
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	file := filepath.Join("views", filename)

	tmpl, err := template.ParseFiles(file)
	if err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}
