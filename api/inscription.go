// ©matthews-crypto
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Définissez une structure pour représenter les données de l'utilisateur
type Users struct {
	Prenom      string `json:"prenom"`
	Nom         string `json:"nom"`
	Telephone   string `json:"telephone"`
	MotDePasse  string `json:"motDePasse"`
	VenteStatut bool   `json:"venteStatut"`
}

// Créez une variable globale pour stocker la collection d'utilisateurs
var personneCollection *mongo.Collection

func inscriptionHandler(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	personneCollection = client.Database("Diayma").Collection("personne")

	// Vérifiez que la méthode de la requête est POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed) // Renvoyez un code d'erreur 405
		w.Write([]byte("Method not allowed"))
		return
	}

	// Créez une variable pour stocker les données de l'utilisateur envoyées par flutter
	var user Users

	// Décodez le corps de la requête en JSON et remplissez la variable user
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Renvoyez un code d'erreur 400
		w.Write([]byte("Invalid data"))
		return
	}

	// Hachez le mot de passe avec le coût par défaut
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.MotDePasse), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	user.MotDePasse = string(hashedPassword)

	personneBSON, err := bson.Marshal(bson.M{
		"login":       user.Telephone,
		"password":    user.MotDePasse,
		"nom":         user.Nom,
		"prenom":      user.Prenom,
		"venteStatut": false,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	insertResultPersonne, err := personneCollection.InsertOne(context.Background(), personneBSON)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	fmt.Println("Inserted a user with ID:", insertResultPersonne.InsertedID, "in personne collection")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User created successfully"))

}

// ©matthews-crypto
