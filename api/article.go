package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type Article struct {
	Titre       string `json:"titre"`
	Description string `json:"description"`
	Prix        int    `json:"prix"`
	Vendeur     string `json:"vendeur"`
}

// ajout d'un article
func addArticle(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "articles"
	collection := client.Database("Diayma").Collection("article")

	// Lecture du corps de la requête
	var article Article
	err = json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		log.Println("Erreur de lecture du corps de la requête:", err)
		http.Error(w, "Erreur de lecture du corps de la requête", http.StatusBadRequest)
		return
	}

	// Insertion de l'article dans la collection
	_, err = collection.InsertOne(context.Background(), article)
	if err != nil {
		log.Println("Erreur lors de l'insertion de l'article dans la base de données:", err)
		http.Error(w, "Erreur lors de l'insertion de l'article dans la base de données", http.StatusInternalServerError)
		return
	}

	// Réponse JSON indiquant que l'article a été ajouté avec succès
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Article ajouté avec succès"})
}

// supression d'un article
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "articles"
	collection := client.Database("Diayma").Collection("article")

	// Obtenir l'ID de l'article à supprimer à partir des paramètres de la requête
	vars := mux.Vars(r)
	articleID := vars["id"]

	// Créer un filtre pour l'ID de l'article
	filter := bson.M{"_id": articleID}

	// Supprimer l'article de la collection
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println("Erreur lors de la suppression de l'article de la base de données:", err)
		http.Error(w, "Erreur lors de la suppression de l'article de la base de données", http.StatusInternalServerError)
		return
	}

	// Réponse JSON indiquant que l'article a été supprimé avec succès
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Article supprimé avec succès"})
}
