package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
)

func authentification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Récupérer les informations de connexion depuis la requête
	r.ParseForm()
	login := r.Form.Get("login")
	motDePasse := r.Form.Get("motDePasse")
	// Se connecter à la base de données MongoDB
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Récupérer la collection "personne"
	collection := client.Database("Diayma").Collection("personne")

	// Définir le filtre de recherche pour le login et le mot de passe
	filter := bson.M{
		"login":      login,
		"motDePasse": motDePasse,
	}

	// Rechercher l'utilisateur correspondant aux informations de connexion
	var personne struct {
		Login      string `json:"login"`
		MotDePasse string `json:"motDePasse"`
		UserType   string `json:"userType"`
	}
	err = collection.FindOne(context.Background(), filter).Decode(&personne)
	fmt.Println(personne.Login)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// L'utilisateur n'a pas été trouvé dans la collection "personne"
			fmt.Println("Identifiants invalides")
			http.Redirect(w, r, "/Diayma/authentification", http.StatusSeeOther)
			return
		}
		log.Fatal(err)
	}

	// Rediriger l'utilisateur en fonction de son type (acheteur ou vendeur)
	if personne.UserType == "acheteur" {
		// Rediriger vers la page de l'acheteur
		http.Redirect(w, r, "/acheteur/ajout", http.StatusSeeOther)
	} else if personne.UserType == "vendeur" {
		// Rediriger vers la page du vendeur
		http.Redirect(w, r, "/vendeur/ajout", http.StatusSeeOther)
	} else {
		// Type d'utilisateur inconnu
		fmt.Println("Type d'utilisateur inconnu")
		http.Redirect(w, r, "/Diayma/authentification", http.StatusSeeOther)
	}
}
