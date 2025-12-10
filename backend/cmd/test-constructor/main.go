package main

import (
	"log"
	"net/http"
	"test-constructor/internal/auth"
	"test-constructor/internal/database"
	"test-constructor/internal/handlers"
	"test-constructor/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const clientURL = "http://localhost:5173"

func main() {
	database.Connect()

	r := mux.NewRouter()

	r.HandleFunc("/register", auth.RegistrationHandler).Methods("POST")
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")

	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	manager := api.PathPrefix("/manager").Subrouter()
	manager.HandleFunc("/tests", handlers.ManagerTestHandler).Methods("GET")
	manager.HandleFunc("/tests/{id}", handlers.DeleteTest).Methods("POST")

	intern := api.PathPrefix("/intern").Subrouter()
	intern.HandleFunc("/tests", handlers.InternAttemptHandler).Methods("GET")

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
