package main

import (
	"context"
	"log"
	"taskflow/internal/config"
	"taskflow/internal/database"
	"taskflow/internal/handlers"
	"taskflow/internal/repository"
	"taskflow/internal/router"
	"taskflow/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPostgresPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	teamService := services.NewTeamService(teamRepo)
	taskService := services.NewTaskService(taskRepo, teamRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userRepo)
	teamHandler := handlers.NewTeamHandler(teamService)
	taskHandler := handlers.NewTaskHandler(taskService)

	r := router.SetupRouter(authHandler, userHandler, teamHandler, taskHandler, cfg.JWTSecret)

	log.Printf("Server starting on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
