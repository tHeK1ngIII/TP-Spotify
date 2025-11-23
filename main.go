package main

import (
	"log"
	"net/http"
	"os"
	"tpspotify/router"
)

func main() {
	// Charge le routeur
	r := router.New()
	// Définition du port (par défaut 8080 si aucune variable d'environnement n'est définie)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Démarrage du serveur sur le port http://localhost:%s\n", port)

	// Démarrage du serveur avec gestion des erreurs
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Erreur lors du démarrage du serveur :", err)
		/*http.Handle("/asset/image/", http.StripPrefix("/asset/image/", http.FileServer(http.Dir("image"))))*/

	}

}
