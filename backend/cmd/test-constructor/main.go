package main

import (
	"log"
	"net/http"
	"test-constructor/internal/database"
	"test-constructor/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const clientURL = "http://localhost:5173"

func main() {
	database.Connect()

	r := mux.NewRouter()

	r.HandleFunc("/register", handlers.Register).Methods("POST")
	r.HandleFunc("/login", handlers.Login).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{clientURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(r)

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
