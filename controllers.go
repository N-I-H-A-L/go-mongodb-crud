package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

// Struct for User
type User struct {
	Name string `json:"name"`
	City string `json:"city"`
	Age  int    `json:"age"`
}

// Get mongo client from db.go file
var userCollection *mongo.Collection = db().Database("mydb").Collection("Users")

// Function to create user profile
func CreateProfile(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Variable to store the incoming user profile
	var person User

	// Decode the incoming JSON request body into the 'person' variable
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		// Handle error if decoding fails
		fmt.Print(err)
		return
	}

	// Insert the person profile into MongoDB
	insertResult, err := userCollection.InsertOne(context.TODO(), person)
	if err != nil {
		// Log error if insertion fails
		log.Fatal(err)
		return
	}

	fmt.Println(insertResult)

	// Return the MongoDB InsertedID as JSON response
	json.NewEncoder(w).Encode(insertResult.InsertedID)
}
