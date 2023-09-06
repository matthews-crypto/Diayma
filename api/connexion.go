// ©matthews-crypto
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Prenom   string `json:"prenom"`
	Nom      string `json:"nom"`
}

type LoginCredentials struct {
	Login    string `json:"login"`
	Password string `json:"motDePasse"`
}

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

	filter := bson.M{"login": credentials.Login}
	fmt.Println("filter: ", filter)
	var user User
	fmt.Println("user: ", user)

	err = collection.FindOne(context.Background(), filter).Decode(&user)
	fmt.Println("user decodé: ", user)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	fmt.Println("password given: ", credentials.Password)
	fmt.Println("password hashed: ", user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		// Le mot de passe est incorrect
		http.Error(w, "mot de passe incorrect", http.StatusUnauthorized)
		fmt.Println(err)
		return
	}
	// user connecter
	http.Redirect(w, r, "/acheteur/ajout", http.StatusAccepted)
	// print(user)
	// print("connected")
}

// ©matthews-crypto
