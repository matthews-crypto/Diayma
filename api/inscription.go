package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Définissez une structure pour représenter les données de l'utilisateur
type Users struct {
	Prenom     string `json:"prenom"`
	Nom        string `json:"nom"`
	Telephone  string `json:"telephone"`
	MotDePasse string `json:"motDePasse"`
}

// Créez une variable globale pour stocker la collection d'utilisateurs
var personneCollection *mongo.Collection
var vendeurCollection *mongo.Collection

func inscriptionHandler(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	personneCollection = client.Database("Diayma").Collection("personne")
	vendeurCollection = client.Database("Diayma").Collection("vendeur")

	// Vérifiez que la méthode de la requête est POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed) // Renvoyez un code d'erreur 405
		w.Write([]byte("Method not allowed"))      // Renvoyez un message d'erreur
		return                                     // Terminez la fonction
	}

	// Créez une variable pour stocker les données de l'utilisateur envoyées par flutter
	var user Users

	// Décodez le corps de la requête en JSON et remplissez la variable user
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Renvoyez un code d'erreur 400
		w.Write([]byte("Invalid data"))      // Renvoyez un message d'erreur
		return                               // Terminez la fonction
	}

	// Convertissez la variable user en deux documents BSON pour les insérer dans les collections correspondantes
	personneBSON, err := bson.Marshal(bson.M{
		"telephone":  user.Telephone,
		"motDePasse": user.MotDePasse,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Renvoyez un code d'erreur 500
		w.Write([]byte("Internal server error"))      // Renvoyez un message d'erreur
		return                                        // Terminez la fonction
	}

	vendeurBSON, err := bson.Marshal(bson.M{
		"telephone": user.Telephone,
		"nom":       user.Nom,
		"prenom":    user.Prenom,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Renvoyez un code d'erreur 500
		w.Write([]byte("Internal server error"))      // Renvoyez un message d'erreur
		return                                        // Terminez la fonction
	}

	// Insérez les documents BSON dans les collections d'utilisateurs
	insertResultPersonne, err := personneCollection.InsertOne(context.Background(), personneBSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Renvoyez un code d'erreur 500
		w.Write([]byte("Internal server error"))      // Renvoyez un message d'erreur
		return                                        // Terminez la fonction
	}

	insertResultVendeur, err := vendeurCollection.InsertOne(context.Background(), vendeurBSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Renvoyez un code d'erreur 500
		w.Write([]byte("Internal server error"))      // Renvoyez un message d'erreur
		return                                        // Terminez la fonction
	}

	fmt.Println("Inserted a user with ID:", insertResultPersonne.InsertedID, "in personne collection")
	fmt.Println("Inserted a user with ID:", insertResultVendeur.InsertedID, "in vendeur collection")

	w.WriteHeader(http.StatusOK)                 // Renvoyez un code de succès 200
	w.Write([]byte("User created successfully")) // Renvoyez un message de succès

}
