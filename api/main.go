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

func enableCors(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "DELETE")
}

func main() {

	router := mux.NewRouter()

	//Inscription ,authentification
	router.HandleFunc("/api/login", login).Methods("POST")
	router.HandleFunc("/api/inscription", inscriptionHandler).Methods("POST")

	// Définition de l'URL  pour les ajouts
	router.HandleFunc("/acheteur/ajout", addAcheteur).Methods("POST")
	router.HandleFunc("/vendeur/ajout", addVendeur).Methods("POST")
	// router.HandleFunc("/artcicle/ajout", addArticle).Methods("POST")

	// Définition de l'URL  pour les suppression
	router.HandleFunc("/acheteur/supprime/{id}", deleteAcheteur).Methods("DELETE")
	router.HandleFunc("/vendeur/supprime/{id}", deleteVendeur).Methods("DELETE")
	// router.HandleFunc("/ article/supprime/{id}", deleteArticle).Methods("DELETE")

	//DEfinition de l'URL pour les modifications
	router.HandleFunc("/acheteur/update/{telephone}v4.", modifierAcheteur).Methods("PUT")

	//Définition de l'URL de lecture
	router.HandleFunc("/acheteur/lire/elementListe/{telephone}", uniqueAcheteur).Methods("GET")
	router.HandleFunc("/acheteur/lire/Liste", listeAcheteur).Methods("GET")
	router.HandleFunc("/vendeur/lire/Liste", listeVendeur).Methods("GET")

	// router pour les articles(groupage)
	api := router.PathPrefix("/api").Subrouter()
	{
		// spécification du groupage des articles
		api.HandleFunc("/products", getProducts).Methods("GET")
		api.HandleFunc("/products", addProduct).Methods("POST")
		api.HandleFunc("/products/{name}", updateProduct).Methods("PUT")
		api.HandleFunc("/products/{name}", deleteProduct).Methods("DELETE")
	}

	//Démarrage du serveur HTTPS
	log.Println("Démarrage du serveur sur le port 192.168.0.70:8080...")
	log.Fatal(http.ListenAndServe("192.168.0.70:8080", router))

}
