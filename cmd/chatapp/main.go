package main

import (
	"github.com/ghenah/chatapp/pkg/database"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	config := database.Config{
		ServerName: "mysql:3306",
		User:       "chatappadmin",
		Password:   "asdfasdf88",
		DB:         "chatapp",
	}

	connectionString := database.GetConnectionString(config)
	err := database.Connect(connectionString)
	if err != nil {
		panic(err.Error())
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/create", createPerson).Methods("POST")
	router.HandleFunc("/get/{id}", getPersonByID).Methods("GET")

	log.Println("Starting the HTTP server on port 8081")
	log.Fatal(http.ListenAndServe(":8081", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var person database.Person
	json.Unmarshal(requestBody, &person)

	fmt.Println(person)

	database.Connector.Create(person)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(person)
}

func getPersonByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	var person database.Person
	database.Connector.First(&person, key)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

// func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
// 	(*w).header().Set("Access-Control-Allow-Origin", "*")
// 	(*w).header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
// 	(*w).header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")

// }