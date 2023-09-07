// Le code du serveur backend en Golang

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Prenom   string `json:"prenom"`
	Nom      string `json:"nom"`
	Token    string `json:"token"`
}

type LoginCredentials struct {
	Login    string `json:"login"`
	Password string `json:"motDePasse"`
}

var secretKey = []byte("secret")

func login(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

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
		http.Error(w, "mot de passe incorrect", http.StatusUnauthorized)
		fmt.Println(err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"prenom":   user.Prenom,
		"nom":      user.Nom,
		"login":    user.Login,
		"password": user.Password,
		"token":    user.Token,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(24 * time.Hour).Unix()})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Erreur de création du token", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"token": tokenString})

	res, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.Login},
		bson.D{
			{"$set", bson.D{{"token", tokenString}}},
		},
	)
	if err != nil {
		log.Println("Erreur de mise à jour du token:", err)
	}
	fmt.Println("Nombre de documents mis à jour:", res.ModifiedCount)
}
