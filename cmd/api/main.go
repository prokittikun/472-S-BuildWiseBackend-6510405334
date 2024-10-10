package main

import (
	"boonkosang/internal/adapters/postgres"
	"boonkosang/internal/adapters/rest"
	"boonkosang/internal/infrastructure/database"
	"boonkosang/internal/infrastructure/server"
	"boonkosang/internal/usecase"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Warning: No .env file found")
	}

	// Now that env vars are loaded, we can use getEnv
	fmt.Println("Boonkosang API", getEnv("DB_HOST", "beer"))

	// Create a new configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "general"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewSQLxDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.CloseSQLxDB(db)

	app := server.NewFiberServer()

	userRepo := postgres.NewUserRepository(db)
	jwtSecret := getEnv("JWT_SECRET", "your_default_secret")
	jwtExpiration := getEnvAsDuration("JWT_EXPIRATION", 15*time.Minute)
	userUseCase := usecase.NewUserUsecase(userRepo, jwtSecret, jwtExpiration)
	UserHandler := rest.NewUserHandler(userUseCase)
	UserHandler.UserRoutes(app)

	projectRepo := postgres.NewProjectRepository(db)
	projectUseCase := usecase.NewProjectUsecase(projectRepo)
	projectHandler := rest.NewProjectHandler(projectUseCase)
	projectHandler.ProjectRoutes(app)

	clientRepo := postgres.NewClientRepository(db)
	clientUseCase := usecase.NewClientUsecase(clientRepo)
	clientHandler := rest.NewClientHandler(clientUseCase)
	clientHandler.ClientRoutes(app)

	materialRepo := postgres.NewMaterialRepository(db)
	materialUseCase := usecase.NewMaterialUsecase(materialRepo)
	materialHandler := rest.NewMaterialHandler(materialUseCase)
	materialHandler.MaterialRoutes(app)

	jobRepo := postgres.NewJobRepository(db)
	jobUseCase := usecase.NewJobUsecase(jobRepo, materialRepo)
	jobHandler := rest.NewJobHandler(jobUseCase)
	jobHandler.JobRoutes(app)

	port := getEnv("PORT", "8004")
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Helper function to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to read an environment variable as an integer or return a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// Helper function to read an environment variable as a duration or return a default value
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
