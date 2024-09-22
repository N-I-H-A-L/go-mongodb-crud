package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// Get the profile of a user via user's name
func getUserProfile(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Variable to store the incoming request body (user information)
	var body User

	// Decode the incoming JSON request body into 'body'
	e := json.NewDecoder(r.Body).Decode(&body)
	if e != nil {
		fmt.Print(e)
		return
	}

	// Variable to store the MongoDB result (as an unordered BSON document)
	var result primitive.M // primitive.M is a map representing BSON document

	// Search MongoDB collection for a user with the given name
	err := userCollection.FindOne(context.TODO(), bson.D{{"name", body.Name}}).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Return the MongoDB result (as a map) in JSON format
	json.NewEncoder(w).Encode(result)
}

func updateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Define a struct for the update body
	type updateBody struct {
		Name string `json:"name"` // Value that has to be matched
		City string `json:"city"` // Value that has to be modified
	}

	var body updateBody

	// Decode the incoming JSON request body into the 'body' variable
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		// Handle error if decoding fails
		fmt.Print(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		//http.StatusBadRequest -> it is status code, set as 400. You can directly write it as 400, no issue.
		return
	}

	// Define the filter to find the document with the specified name
	filter := bson.D{{Key: "name", Value: body.Name}}

	// Define the update operation to set the new city value
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "city", Value: body.City}}}}

	// Define options for the FindOneAndUpdate operation
	returnOpt := options.FindOneAndUpdate().SetReturnDocument(options.After)
	// returnOpt is configured to return the updated document by setting ReturnDocument: options.After.

	// Perform the update operation
	updateResult := userCollection.FindOneAndUpdate(context.TODO(), filter, update, returnOpt)

	// Check for errors and decode the result
	var result primitive.M
	err = updateResult.Decode(&result)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	// Return the updated document as JSON response
	json.NewEncoder(w).Encode(result)
}

func deleteProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get parameter value as string
	params := mux.Vars(r)["id"]

	// Convert params to MongoDB Hex ID
	_id, err := primitive.ObjectIDFromHex(params)
	if err != nil {
		fmt.Print(err.Error())
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Specify options for deletion, if needed
	opts := options.Delete().SetCollation(&options.Collation{
		Locale: "en", // Example locale, specify as needed
	})

	// Perform the delete operation
	res, err := userCollection.DeleteOne(context.TODO(), bson.D{{Key: "_id", Value: _id}}, opts)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}

	// Return the number of documents deleted
	json.NewEncoder(w).Encode(res.DeletedCount)
}
