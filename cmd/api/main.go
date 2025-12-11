package main

import (
	"log"
	"net/http"
	"github.com/MH-Cognition/mhc-backend-lms-identity-service.git/internal/repository/inmemory"
	httptransport "github.com/MH-Cognition/mhc-backend-lms-identity-service.git/internal/transport/http"
	"github.com/MH-Cognition/mhc-backend-lms-identity-service.git/internal/usecase"
)

func main() {

	userRepo := inmemory.NewInMemoryUserRepository()

	userUC := usecase.NewUserUsecase(userRepo)

	userHandler := httptransport.NewUserHandler(userUC)

	r := httptransport.NewRouter(userHandler)

	log.Println("🚀 Server running on :8080")
	http.ListenAndServe(":8080", r)

}