package main

import (
	"log"
	"net/http"

	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/repository/inmemory"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/transport/httptransport"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/usecase"
)

func main() {

	userRepo := inmemory.NewInMemoryUserRepository()

	userUC := usecase.NewUserUsecase(userRepo)

	userHandler := httptransport.NewUserHandler(userUC)

	r := httptransport.NewRouter(userHandler)

	log.Println("🚀 Server running on :8080")
	http.ListenAndServe(":8080", r)

}
