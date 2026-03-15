package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	v1 "practice7/internal/controller/http/v1"
	"practice7/internal/usecase"
	"practice7/internal/usecase/repo"
	"practice7/pkg/postgres"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	pg, err := postgres.New()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repo.NewUserRepo(pg)
	userUseCase := usecase.NewUserUseCase(userRepo)

	r := gin.Default()

	api := r.Group("/api")
	v1.NewUserRoutes(api, userUseCase)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}

}
