package controller

import (
	"html/template"
	"net/http"
)

// Fonction utilitaire pour afficher un template
func renderTemplate(w http.ResponseWriter, filename string, data map[string]string) {
	tmpl := template.Must(template.ParseFiles("template/" + filename))
	tmpl.Execute(w, data)
}

// Page d'accueil
func Home(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title":   "Accueil",
		"Message": "Clique sur une image pour voir les infos Spotify 🎵",
	}
	renderTemplate(w, "index.html", data)
}

// Page Damso
func DamsoAlbum(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title":   "Damso",
		"Message": "Voici la page dédiée à Damso 🎤",
	}
	renderTemplate(w, "damso.html", data)
}

// Page Laylow
func LaylowSong(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Title":   "Laylow",
		"Message": "Voici la page dédiée à Laylow 🎧",
	}
	renderTemplate(w, "laylow.html", data)
}
