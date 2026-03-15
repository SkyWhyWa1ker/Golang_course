package main

import (
	"log"
	"net/http"

	"practice5/internal/config"
	"practice5/internal/handler"
	"practice5/internal/middleware"
	"practice5/internal/repository"
	"practice5/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", userHandler.Health)
	mux.HandleFunc("/users", userHandler.Users)
	mux.HandleFunc("/users/", userHandler.UserByID)
	mux.HandleFunc("/common-friends", userHandler.CommonFriends)

	finalHandler := middleware.LoggingMiddleware(
		middleware.APIKeyMiddleware(cfg.APIKey, mux),
	)

	log.Println("Starting the Server on :8080")
	if err := http.ListenAndServe(":8080", finalHandler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
