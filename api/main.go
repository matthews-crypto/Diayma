package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// donction de connexion a la base de donnees
func connectToDatabase() (*mongo.Client, error) {
	// Configuration de la connexion MongoDB
	uri := "mongodb+srv://matthews:Passer2023@diayma.zibf8gu.mongodb.net/?retryWrites=true&w=majority"
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Établir une connexion à la base de données
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}
func main() {

	router := mux.NewRouter()

	//partie inscription
	router.HandleFunc("/Diayma/inscription", inscription).Methods("POST")

	//partie connexion
	// router.HandleFunc("/authentification", handleLogin).Methods("POST")

	// Définition de l'URL  pour les ajouts
	router.HandleFunc("/acheteur/ajout", addAcheteur).Methods("POST")
	router.HandleFunc("/vendeur/ajout", addVendeur).Methods("POST")
	router.HandleFunc("/artcicle/ajout", addArticle).Methods("POST")

	// Définition de l'URL  pour les suppression
	router.HandleFunc("/acheteur/supprime/{id}", deleteAcheteur).Methods("DELETE")
	router.HandleFunc("/vendeur/supprime/{id}", deleteVendeur).Methods("DELETE")
	router.HandleFunc("/ article/supprime/{id}", deleteArticle).Methods("DELETE")

	//DEfinition de l'URL pour les modifications
	router.HandleFunc("/acheteur/update/{telephone}", modifierAcheteur).Methods("PUT")

	//Définition de l'URL de lecture
	router.HandleFunc("/acheteur/lire/elementListe/{telephone}", uniqueAcheteur).Methods("GET")
	router.HandleFunc("/acheteur/lire/Liste", listeAcheteur).Methods("GET")

	// Démarrage du serveur HTTP
	log.Println("Démarrage du serveur sur le port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
