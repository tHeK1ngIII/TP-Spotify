package controller

import (
	"html/template"
	"net/http"
)

// --- Templates ---
func renderTemplate(w http.ResponseWriter, filename string, data interface{}) {
	tmpl := template.Must(template.ParseFiles("template/" + filename))
	_ = tmpl.Execute(w, data)
}

// --- Home ---
func Home(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", nil)
}

// --- Laylow ---
func LaylowSong(w http.ResponseWriter, r *http.Request) {
	// Exemple minimal
	renderTemplate(w, "laylow.html", nil)
}
