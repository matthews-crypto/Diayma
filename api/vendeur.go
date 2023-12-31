package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type Vendeur struct {
	Nom           string   `json:"nom"`
	Prenom        string   `json:"prenom"`
	Login         string   `json:"login"`
	Telephone     int      `json:"telephone"`
	MoyenPaiement []string `json:"moyenPaiement"`
	NomBoutique   []string `json:"nomBoutique"`
	Cin           int      `json:"cin"`
	Email         string   `json:"email"`
	TypeVendeur   []string `json:"typeVendeur"`
}

// ajout d'un vendeur
func addVendeur(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "vendeurs"
	collection := client.Database("Diayma").Collection("vendeur")

	// Lecture du corps de la requête
	var vendeur Vendeur
	err = json.NewDecoder(r.Body).Decode(&vendeur)
	if err != nil {
		log.Println("Erreur de lecture du corps de la requête:", err)
		http.Error(w, "Erreur de lecture du corps de la requête", http.StatusBadRequest)
		return
	}

	// Insertion de l'vendeur dans la collection
	_, err = collection.InsertOne(context.Background(), vendeur)
	if err != nil {
		log.Println("Erreur lors de l'insertion de l'vendeur dans la base de données:", err)
		http.Error(w, "Erreur lors de l'insertion de l'vendeur dans la base de données", http.StatusInternalServerError)
		return
	}

	// Réponse JSON indiquant que l'vendeur a été ajouté avec succès
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Vendeur ajouté avec succès"})
}

// supression d'un vendeur
func deleteVendeur(w http.ResponseWriter, r *http.Request) {
	enableCors(w, r)
	// Connexion à la base de données
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "vendeurs"
	collection := client.Database("Diayma").Collection("vendeur")

	// Obtenir l'ID de l'vendeur à supprimer à partir des paramètres de la requête
	vars := mux.Vars(r)
	vendeurID := vars["id"]
	telephoneI, err := strconv.Atoi(vendeurID)
	if err != nil {
		fmt.Println("erreur convertion")
	}

	// Créer un filtre pour l'ID de l'vendeur
	filter := bson.M{"telephone": telephoneI}

	// Supprimer l'vendeur de la collection
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println("Erreur lors de la suppression de l'vendeur de la base de données:", err)
		http.Error(w, "Erreur lors de la suppression de l'vendeur de la base de données", http.StatusInternalServerError)
		return
	}

	// Réponse JSON indiquant que l'vendeur a été supprimé avec succès
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Vendeur supprimé avec succès"})
}

func listeVendeur(w http.ResponseWriter, r *http.Request) {
	enableCors(w, r)
	//Connexion à la base de données
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "vendeur"
	collection := client.Database("Diayma").Collection("vendeur")

	// Rechercher tous les documents dans la collection
	cur, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())

	// Créer une liste de vendeurs
	var vendeurs []Vendeur

	// Parcourir les résultats et ajouter chaque document à la liste des vendeurs
	for cur.Next(context.Background()) {
		var vendeur Vendeur
		err := cur.Decode(&vendeur)
		if err != nil {
			log.Fatal(err)
		}
		vendeurs = append(vendeurs, vendeur)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Convertir la liste des vendeurs en JSON
	vendeursJSON, err := json.Marshal(vendeurs)
	if err != nil {
		log.Fatal(err)
	}

	// Définir le type de contenu de la réponse
	w.Header().Set("Content-Type", "application/json")

	// Envoyer la liste des vendeurs en tant que réponse
	_, err = w.Write(vendeursJSON)
	if err != nil {
		log.Fatal(err)
	}
}
