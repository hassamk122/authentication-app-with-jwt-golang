package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hassamk122/authentication-app-with-jwt-golang/config"
	transaction "github.com/hassamk122/authentication-app-with-jwt-golang/internals/Transaction"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/handlers"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/middlewares"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/repo"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/routes"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/services"
	"github.com/hassamk122/authentication-app-with-jwt-golang/internals/store"
)

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	db := config.ConnectToDB(configuration.DatabaseUrl)
	defer db.Close()

	queries := store.New(db)

	userRepo := repo.NewUserRepo(queries)

	verificationCodeRepo := repo.NewVerificationCodeRepo(queries)

	userSessionRepo := repo.NewUserSessionRepo(queries)

	txManager := transaction.NewTxManager[any](db)

	userService := services.NewUserService(*txManager, userRepo, verificationCodeRepo, userSessionRepo)

	handler := handlers.NewHandler(userService)

	mux := http.NewServeMux()

	routes.SetupRoutes(mux, handler)

	serverAddr := fmt.Sprintf(":%s", configuration.ServerPort)

	loggingMux := middlewares.LoggingMiddleware(mux)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: loggingMux,
	}

	log.Printf("Starting server on Port %s", serverAddr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed. Reason: %v", err)
	}
}
