package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	route := mux.NewRouter()                  //Gorilla Mux
	s := route.PathPrefix("/api").Subrouter() //Base Path

	//Routes
	fmt.Println("here")
	s.HandleFunc("/createProfile", CreateProfile).Methods("POST")
	// s.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")
	// s.HandleFunc("/getUserProfile", getUserProfile).Methods("POST")
	// s.HandleFunc("/updateProfile", updateProfile).Methods("PUT")
	// s.HandleFunc("/deleteProfile/{id}", deleteProfile).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", s)) // Run Server
}
