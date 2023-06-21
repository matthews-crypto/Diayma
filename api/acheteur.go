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

type Acheteur struct {
	Nom         string   `json:"nom"`
	Prenom      string   `json:"prenom"`
	Panier      []string `json:"panier"`
	Commandes   []string `json:"commandes"`
	Preferences []string `json:"preferences"`
	Telephone   int      `json:"telephone"`
}

// ajout d'un acheteur
func addAcheteur(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "acheteurs"
	collection := client.Database("Diayma").Collection("acheteur")

	// Lecture du corps de la requête
	var acheteur Acheteur
	err = json.NewDecoder(r.Body).Decode(&acheteur)
	if err != nil {
		log.Println("Erreur de lecture du corps de la requête:", err)
		http.Error(w, "Erreur de lecture du corps de la requête", http.StatusBadRequest)
		return
	}

	// Insertion de l'acheteur dans la collection
	_, err = collection.InsertOne(context.Background(), acheteur)
	if err != nil {
		log.Println("Erreur lors de l'insertion de l'acheteur dans la base de données:", err)
		http.Error(w, "Erreur lors de l'insertion de l'acheteur dans la base de données", http.StatusInternalServerError)
		return
	}

	// Réponse JSON indiquant que l'acheteur a été ajouté avec succès
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Acheteur ajouté avec succès"})
}

// supression d'un acheteur
func deleteAcheteur(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "acheteurs"
	collection := client.Database("Diayma").Collection("acheteur")

	// Obtenir l'ID de l'acheteur à supprimer à partir des paramètres de la requête
	vars := mux.Vars(r)
	telephone := vars["id"]

	// Créer un filtre pour l'ID de l'acheteur
	filter := bson.M{"_id": telephone}

	// Supprimer l'acheteur de la collection
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println("Erreur lors de la suppression de l'acheteur de la base de données:", err)
		http.Error(w, "Erreur lors de la suppression de l'acheteur de la base de données", http.StatusInternalServerError)
		return
	}

	// Réponse JSON indiquant que l'acheteur a été supprimé avec succès
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Acheteur supprimé avec succès"})
}

// Modification.
func modifierAcheteur(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de requête à l'aide de mux
	params := mux.Vars(r)
	telephone := params["telephone"]

	// Connexion à la base de données MongoDB
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "acheteurs"
	collection := client.Database("Diayma").Collection("acheteur")

	// Créer un filtre pour l'ID de l'acheteur
	telephoneI, err := strconv.Atoi(telephone)
	if err != nil {
		fmt.Println("erreur convertion")
	}

	filter := bson.M{"telephone": telephoneI}

	// Décoder les données de la requête JSON pour obtenir les champs mis à jour
	var updatedFields bson.M
	err = json.NewDecoder(r.Body).Decode(&updatedFields)
	if err != nil {
		// Erreur de décodage JSON
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Créer une mise à jour pour les champs spécifiés
	update := bson.M{"$set": updatedFields}

	// Effectuer la mise à jour de l'acheteur dans la base de données
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		// Erreur lors de la mise à jour de l'acheteur
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Répondre avec un message de succès
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Acheteur mis à jour avec succès")
}
func uniqueAcheteur(w http.ResponseWriter, r *http.Request) {
	// Récupérer les paramètres de requête à l'aide de mux
	params := mux.Vars(r)
	telephone := params["telephone"]

	// Connexion à la base de données MongoDB
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "acheteurs"
	collection := client.Database("Diayma").Collection("acheteur")

	// Créer un filtre pour l'ID de l'acheteur
	telephoneI, err := strconv.Atoi(telephone)
	if err != nil {
		fmt.Println("erreur convertion")
	}
	filter := bson.M{"telephone": telephoneI}

	// Rechercher l'acheteur dans la base de données
	var acheteur Acheteur
	err = collection.FindOne(context.Background(), filter).Decode(&acheteur)
	if err != nil {
		// Acheteur non trouvé
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Convertir l'acheteur en JSON
	response, err := json.Marshal(acheteur)
	if err != nil {
		// Erreur lors de la conversion en JSON
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Répondre avec l'acheteur en tant que JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func listeAcheteur(w http.ResponseWriter, r *http.Request) {
	// Connexion à la base de données MongoDB
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "acheteurs"
	collection := client.Database("Diayma").Collection("acheteur")

	// Définition du contexte
	ctx := context.Background()

	// Définition du filtre (ici, vide pour récupérer tous les acheteurs)
	filter := bson.D{}

	// Définition d'un slice pour stocker les acheteurs
	var acheteurs []Acheteur

	// Exécution de la requête de recherche
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Décodage des résultats dans le slice d'acheteurs
	if err := cursor.All(ctx, &acheteurs); err != nil {
		log.Fatal(err)
	}

	// Parcours des acheteurs et affichage des informations
	for _, acheteur := range acheteurs {
		fmt.Println("Nom:", acheteur.Nom)
		fmt.Println("Prénom:", acheteur.Prenom)
		fmt.Println("Panier:", acheteur.Panier)
		fmt.Println("Commandes:", acheteur.Commandes)
		fmt.Println("Préférences:", acheteur.Preferences)
		fmt.Println("Téléphone:", acheteur.Telephone)
		fmt.Println("----------------------------")
	}
}
