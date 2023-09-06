// // Le code du serveur backend en Golang

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// La structure qui représente un produit
type Product struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
	Image string `json:"image"`
}

// // La fonction qui se connecte à la base de données MongoDB
// func connectDB() *mongo.Collection {
// 	// Créer un contexte avec un timeout de 10 secondes
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// // Se connecter au serveur MongoDB local
// 	// client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }

// 	// // Retourner la collection "products" de la base de données "flutter_db"
// 	// return client.Database("flutter_db").Collection("products")
// }

// La fonction qui récupère tous les produits depuis la base de données
func getProducts(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "vendeurs"
	collection := client.Database("Diayma").Collection("article")

	// Créer un slice de produits
	var products []Product

	// Trouver tous les documents dans la collection et les stocker dans le slice
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var product Product
		cursor.Decode(&product)
		products = append(products, product)
	}

	// Encoder le slice de produits au format JSON et l'écrire dans la réponse avec un code 200
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// La fonction qui ajoute un produit dans la base de données
func addProduct(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "vendeurs"
	collection := client.Database("Diayma").Collection("article")

	// Créer un objet produit à partir du corps de la requête
	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	// Insérer le produit dans la collection
	res, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		log.Fatal(err)
	}

	// Encoder l'identifiant du produit inséré au format JSON et l'écrire dans la réponse avec un code 201
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"_id": res.InsertedID})
}

// La fonction qui modifie un produit dans la base de données
func updateProduct(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "vendeurs"
	collection := client.Database("Diayma").Collection("article")

	// Récupérer le nom du produit à modifier depuis le paramètre d'URL
	vars := mux.Vars(r)
	name := vars["name"]

	// Créer un objet produit à partir du corps de la requête
	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	// Mettre à jour le produit dans la collection
	res, err := collection.UpdateOne(
		context.Background(),
		bson.M{"name": name},
		bson.D{
			{"$set", bson.D{{"price", product.Price}, {"image", product.Image}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Encoder le nombre de documents modifiés au format JSON et l'écrire dans la réponse avec un code 200
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"modified_count": res.ModifiedCount})
}

// La fonction qui supprime un produit dans la base de données
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	// Sélectionner la collection "vendeurs"
	collection := client.Database("Diayma").Collection("article")

	// Récupérer le nom du produit à supprimer depuis le paramètre d'URL
	vars := mux.Vars(r)
	name := vars["name"]

	// Supprimer le produit dans la collection
	res, err := collection.DeleteOne(context.Background(), bson.M{"name": name})
	if err != nil {
		log.Fatal(err)
	}

	// Encoder le nombre de documents supprimés au format JSON et l'écrire dans la réponse avec un code 200
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"deleted_count": res.DeletedCount})
}
