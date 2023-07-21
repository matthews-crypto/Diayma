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

//	func enableCors(w http.ResponseWriter, r *http.Request) {
//		(w).Header().Set("Access-Control-Allow-Origin", "*")
//		(w).Header().Set("Access-Control-Allow-Methods", "DELETE")
//		if r.Method == "OPTIONS" {
//			// Répondre avec un en-tête CORS réussi pour les requêtes OPTIONS
//			(w).Header().Set("Access-Control-Allow-Headers", "Content-Type") // Vous pouvez ajouter d'autres en-têtes si nécessaire
//			return
//		}
//	}
func gestionCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Autoriser les requêtes provenant de n'importe quel domaine (remplacez '*' par le domaine spécifique si nécessaire)
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Autoriser les méthodes HTTP spécifiées (dans ce cas, nous utilisons DELETE)
		w.Header().Set("Access-Control-Allow-Methods", "DELETE")

		// Vérifier si la requête est une requête OPTIONS
		if r.Method == "OPTIONS" {
			// Répondre avec un en-tête CORS réussi pour les requêtes OPTIONS
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Vous pouvez ajouter d'autres en-têtes si nécessaire
			return
		}

		// Appeler le gestionnaire HTTP réel
		next(w, r)
	}
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
	router.HandleFunc("/vendeur/supprime/{id}", gestionCors(deleteVendeur)).Methods("DELETE")
	router.HandleFunc("/ article/supprime/{id}", deleteArticle).Methods("DELETE")

	//DEfinition de l'URL pour les modifications
	router.HandleFunc("/acheteur/update/{telephone}v4.", modifierAcheteur).Methods("PUT")

	//Définition de l'URL de lecture
	router.HandleFunc("/acheteur/lire/elementListe/{telephone}", uniqueAcheteur).Methods("GET")
	router.HandleFunc("/acheteur/lire/Liste", listeAcheteur).Methods("GET")
	router.HandleFunc("/vendeur/lire/Liste", gestionCors(listeVendeur)).Methods("GET")

	// Démarrage du serveur HTTP
	log.Println("Démarrage du serveur sur le port 8080...")
	log.Fatal(http.ListenAndServe("192.168.0.82:8080", router))
}
