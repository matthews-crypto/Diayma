package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

func inscription(w http.ResponseWriter, r *http.Request) {
	// Récupérer les informations de l'utilisateur depuis la requête
	login := r.FormValue("login")
	password := r.FormValue("password")
	typeUser := r.FormValue("type user")
	nom := r.FormValue("nom")
	prenom := r.FormValue("prenom")
	telephone := r.FormValue("telephone")

	// Se connecter à la base de données MongoDB
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Récupérer la collection "personne"
	personneCollection := client.Database("Diayma").Collection("personne")

	// Vérifier si l'utilisateur existe déjà dans la collection "personne"
	filter := bson.M{"login": login}
	var existingUser bson.M
	err = personneCollection.FindOne(context.TODO(), filter).Decode(&existingUser)
	if err == nil {
		// L'utilisateur existe déjà
		fmt.Println("L'utilisateur existe déjà")
		http.Redirect(w, r, "/inscription", http.StatusSeeOther)
		return
	} else if err != mongo.ErrNoDocuments {
		// Erreur lors de la recherche de l'utilisateur
		log.Fatal(err)
	}

	// Insérer l'utilisateur dans la collection "personne"
	newUser := bson.M{
		"login":      login,
		"motDePasse": password,
		"typeUser":   typeUser,
	}
	_, err = personneCollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		log.Fatal(err)
	}

	// Récupérer l'ID de l'utilisateur nouvellement inscrit
	var insertedUser struct {
		ID primitive.ObjectID `bson:"_id"`
	}
	err = personneCollection.FindOne(context.TODO(), filter).Decode(&insertedUser)
	if err != nil {
		log.Fatal(err)
	}

	// Récupérer la collection "acheteur"
	acheteurCollection := client.Database("Diayma").Collection("acheteur")

	// Insérer les informations personnelles de l'acheteur dans la collection "acheteur"
	newAcheteur := bson.M{
		"personneID": insertedUser.ID,
		"nom":        nom,
		"prenom":     prenom,
		"telephone":  telephone,
	}
	_, err = acheteurCollection.InsertOne(context.TODO(), newAcheteur)
	if err != nil {
		log.Fatal(err)
	}

	// Rediriger vers la page de connexion après l'inscription réussie
	http.Redirect(w, r, "/authentification", http.StatusSeeOther)
}

// authentification
func handleLogin(w http.ResponseWriter, r *http.Request) {
	// Récupérer les informations de connexion depuis la requête
	login := r.FormValue("login")
	password := r.FormValue("password")

	// Se connecter à la base de données MongoDB
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Récupérer la collection "personne"
	personneCollection := client.Database("Diayma").Collection("personne")

	// Définir le filtre de recherche pour le login et le mot de passe
	filter := bson.M{
		"login":        login,
		"motDePasse":   password,
		"typePersonne": "acheteur",
	}

	// Rechercher l'utilisateur correspondant aux informations de connexion
	var personne bson.M
	err = personneCollection.FindOne(context.TODO(), filter).Decode(&personne)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// L'utilisateur n'a pas été trouvé dans la collection "personne"
			fmt.Println("Identifiants invalides")
			http.Redirect(w, r, "/authentification", http.StatusSeeOther)
			return
		}
		log.Fatal(err)
	}

	// Récupérer l'ID de l'utilisateur
	userID, ok := personne["_id"].(primitive.ObjectID)
	if !ok {
		fmt.Println("ID d'utilisateur invalide")
		http.Redirect(w, r, "/authentification", http.StatusSeeOther)
		return
	}

	// Rediriger vers la page de l'acheteur en incluant l'ID de l'utilisateur
	http.Redirect(w, r, "/acheteur?id="+userID.Hex(), http.StatusSeeOther)
}
