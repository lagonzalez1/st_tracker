package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"tracker/app/database"
	"tracker/app/services"
	"tracker/app/transport"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	db, db_err := database.ConnectDB()
	if db_err != nil {
		log.Fatalf("Database connection failed: %v", db_err)
	}
	defer db.Close()

	r := mux.NewRouter()

	authService := services.NewAuthService(db)
	authHandler := transport.NewAuthHandler(authService)

	r.HandleFunc("/hello", hello).Methods("GET")
	r.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")

	handler := cors.Default().Handler(r)
	err := http.ListenAndServe(":3333", handler)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
