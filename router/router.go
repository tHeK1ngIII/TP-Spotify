package router

import (
	"net/http"
	"tpspotify/controller"
)

func New() http.Handler {
	mux := http.NewServeMux()

	// Fichiers statiques (CSS + images)
	mux.Handle("/style.css", http.FileServer(http.Dir("web")))
	mux.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("web/image"))))

	// Routes
	mux.HandleFunc("/", controller.Home)
	mux.HandleFunc("/album/damso", controller.DamsoAlbum)
	mux.HandleFunc("/track/laylow", controller.LaylowSong)

	return mux
}
