package main

import (
	"log"
	"net/http"
	"tpspotify/controller"
)

func main() {
	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/image/", http.StripPrefix("/image/", http.FileServer(http.Dir("image"))))

	// Routes
	http.HandleFunc("/", controller.Home)
	http.HandleFunc("/album/damso", controller.DamsoAlbum)
	http.HandleFunc("/track/laylow", controller.LaylowSong)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
