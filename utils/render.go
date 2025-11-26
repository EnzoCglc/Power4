package utils

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Render parses and executes an HTML template with the provided data.
func Render(w http.ResponseWriter, filename string, data interface{}) {
	// Construct full path to template file in views directory
	file := filepath.Join("views", filename)

	// Parse the template file
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	// Execute the template with provided data and write to response
	if err := tmpl.Execute(w, data); err != nil {
		log.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}
