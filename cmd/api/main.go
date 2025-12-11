package main

import (
	"log"
	"net/http"

	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/repository/inmemory"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/transport/httptransport"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/usecase"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/repository/postgres"
)

func main() {

	pool, err := postgres.Connect()
	if err != nil {
		log.Fatal("❌ Failed to connect to Postgres:", err)
	}

	log.Println("✅ Connected to Postgres!",pool)

	userRepo := inmemory.NewInMemoryUserRepository()

	userUC := usecase.NewUserUsecase(userRepo)

	userHandler := httptransport.NewUserHandler(userUC)

	r := httptransport.NewRouter(userHandler)

	log.Println("🚀 Server running on :8080")
	http.ListenAndServe(":8080", r)

}
