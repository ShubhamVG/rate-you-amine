package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var db = &Database{dbName: "database.db"}

func init() {
	if err := db.init(); err != nil {
		log.Fatalln("Failed to init db")
	}
}

func main() {
	defer db.close()

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/ping", pingHandler)
	router.HandleFunc("/api/v1/login", loginHandler).Methods("POST")
	router.HandleFunc("/api/v1/signup", signUpHandler)
	router.HandleFunc("/api/v1/tier/{identifier}", tierHandler) // TODO: name
	router.HandleFunc("/api/v1/create-tier", createTierHandler)
	router.HandleFunc("/api/v1/delete-tier", deleteTierHandler)

	log.Fatalln(http.ListenAndServe(":8080", router))
}
