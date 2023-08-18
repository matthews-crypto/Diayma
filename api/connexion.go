package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"motDePasse"`
	Type     string `json:"userType"`
}

type LoginCredentials struct {
	Login    string `json:"login"`
	Password string `json:"motDePasse"`
	Type     string `json:"userType"`
}

// func inscription(w http.ResponseWriter, r *http.Request) {
// 	// Your existing registration code
// }

func login(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "vendeurs"
	collection := client.Database("Diayma").Collection("personne")
	var credentials LoginCredentials
	err = json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	filter := bson.M{"login": credentials.Login, "motDePasse": credentials.Password, "userType": credentials.Type}
	fmt.Println(filter)
	var user User
	err = collection.FindOne(context.Background(), filter).Decode(&user)
	fmt.Println(user)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	} else {
		http.Error(w, "Connexte avec succes", http.StatusAccepted)
	}
	// fmt.Println(user.Password)

	// if user.Password == credentials.Password {
	// 	if user.Type == "acheteur" {
	// 		http.Redirect(w, r, "/acheteur/ajout", http.StatusFound)
	// 	} else if user.Type == "vendeur" {
	// 		http.Redirect(w, r, "/vendeur/ajout", http.StatusFound)
	// 	} else {
	// 		http.Error(w, "Invalid user type", http.StatusBadRequest)
	// 	}
	// } else {
	// 	http.Error(w, "Login ou mot de passe incorrect", http.StatusUnauthorized)
	// }
}
