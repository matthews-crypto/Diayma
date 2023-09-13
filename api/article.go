// ©matthews-crypto

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type Product struct {
	Nom         string `json:"nom"`
	Quantite    string `json:"quantite"`
	Prix        int    `json:"prix"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())
	collection := client.Database("Diayma").Collection("article")

	var products []Product

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func addProduct(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("Diayma").Collection("article")

	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	res, err := collection.InsertOne(context.Background(), product)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"_id": res.InsertedID})
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("Diayma").Collection("article")

	vars := mux.Vars(r)
	name := vars["name"]

	var product Product
	json.NewDecoder(r.Body).Decode(&product)

	res, err := collection.UpdateOne(
		context.Background(),
		bson.M{"name": name},
		bson.D{
			{"$set", bson.D{{"prix", product.Prix}, {"image", product.Image}}},
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"modified_count": res.ModifiedCount})
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	client, err := connectToDatabase()
	if err != nil {
		log.Println("Erreur de connexion à la base de données:", err)
		http.Error(w, "Erreur de connexion à la base de données", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("Diayma").Collection("article")

	vars := mux.Vars(r)
	name := vars["name"]

	res, err := collection.DeleteOne(context.Background(), bson.M{"name": name})
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"deleted_count": res.DeletedCount})
}

// ©matthews-crypto
