// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/repository/inmemory"
// 	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/repository/postgres"
// 	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/transport/httptransport"
// 	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/usecase"
// )

// func main() {

// 	pool, err := postgres.Connect()
// 	if err != nil {
// 		log.Fatal("❌ Failed to connect to Postgres:", err)
// 	}

// 	log.Println("✅ Connected to Postgres!", pool)

// 	userRepo := inmemory.NewInMemoryUserRepository()

// 	userUC := usecase.NewUserUsecase(userRepo)

// 	userHandler := httptransport.NewUserHandler(userUC)

// 	r := httptransport.NewRouter(userHandler)

// 	log.Println("🚀 Server running on :8080")
// 	http.ListenAndServe(":8080", r)

// }
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/observability/logger"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/observability/tracing"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/repository/inmemory"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/repository/postgres"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/transport/httptransport"
	"github.com/MH-Cognition/golang-microservice-template-clean-architecture/internal/usecase"
)

func main() {
	// ------------------------------------------------
	// Environment
	// ------------------------------------------------
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	serviceName := "user-service"

	// ------------------------------------------------
	// Logger (first)
	// ------------------------------------------------
	appLogger := logger.New(env)

	// ------------------------------------------------
	// Tracing
	// ------------------------------------------------
	shutdown, err := tracing.InitTracer(context.Background(), serviceName)
	if err != nil {
		appLogger.Error(context.Background(), "failed to init tracing", "error", err)
		os.Exit(1)
	}
	defer shutdown(context.Background())

	// ------------------------------------------------
	// Database
	// ------------------------------------------------
	pool, err := postgres.Connect()
	if err != nil {
		appLogger.Error(context.Background(), "failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	appLogger.Info(context.Background(), "postgres connected")

	// ------------------------------------------------
	// Repository
	// ------------------------------------------------
	userRepo := inmemory.NewInMemoryUserRepository()
	_ = pool

	// ------------------------------------------------
	// Usecase
	// ------------------------------------------------
	userUC := usecase.NewUserUsecase(userRepo)

	// ------------------------------------------------
	// Handler
	// ------------------------------------------------
	userHandler := httptransport.NewUserHandler(userUC)

	// ------------------------------------------------
	// Router
	// ------------------------------------------------
	r := httptransport.NewRouter(userHandler, appLogger, env, serviceName)

	// ------------------------------------------------
	// Start server
	// ------------------------------------------------
	appLogger.Info(context.Background(), "server starting", "port", 8080)
	if err := http.ListenAndServe(":8080", r); err != nil {
		appLogger.Error(context.Background(), "server failed", "error", err)
	}
}
