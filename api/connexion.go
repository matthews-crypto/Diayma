package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Type     string `json:"type"` //
}

func inscription(w http.ResponseWriter, r *http.Request) {
	// Parse the request body into a User struct
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Connect to MongoDB
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "acheteurs"
	collectionPersonne := client.Database("Diayma").Collection("personne")

	// Insert the user data into the appropriate collection based on user type
	if user.Type == "acheteur" {
		personne := bson.M{
			"login":    user.Login,
			"password": user.Password,
			"type":     "acheteur",
		}
		_, err = collectionPersonne.InsertOne(context.Background(), personne)
		http.Redirect(w, r, "/acheteur/ajout", http.StatusFound)
	} else if user.Type == "vendeur" {
		personne := bson.M{
			"login":    user.Login,
			"password": user.Password,
			"type":     "vendeur",
		}
		_, err = collectionPersonne.InsertOne(context.Background(), personne)
		http.Redirect(w, r, "/vendeur/ajout", http.StatusFound)

	} else {
		http.Error(w, "Invalid user type", http.StatusBadRequest)
		return
	}

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Return success message
	// w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User registered successfully")
}
