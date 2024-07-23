package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const (
	HOST_ADDR = ":8080"
	DB_NAME   = "database.db"
)

var (
	db     *sql.DB
	router *mux.Router
)

func init() {
	var err error

	if db, err = sql.Open("sqlite3", DB_NAME); err != nil {
		log.Fatalln("Failed to open database")
	}

	router = mux.NewRouter()
	router.HandleFunc("/api/v1/ping", pingHandler).Methods("GET")
	router.HandleFunc("/api/v1/login", loginHandler).Methods("POST")
	router.HandleFunc("/api/v1/signup", signUpHandler).Methods("POST")
	router.HandleFunc("/api/v1/tier/{identifier}", tierHandler).Methods("GET")
	router.HandleFunc("/api/v1/create-tier", createTierHandler).Methods("POST")
	router.HandleFunc("/api/v1/delete-tier", deleteTierHandler).Methods("DELETE")
}

func main() {
	defer db.Close()

	srvr := &http.Server{
		Addr:         HOST_ADDR,
		Handler:      router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	log.Fatalln(srvr.ListenAndServe())
}
