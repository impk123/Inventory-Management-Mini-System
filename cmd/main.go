package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/impk123/Inventory-Management-Mini-System/internal/config"
	"github.com/impk123/Inventory-Management-Mini-System/internal/handlers"
	"github.com/impk123/Inventory-Management-Mini-System/pkg/cache"
	"github.com/impk123/Inventory-Management-Mini-System/pkg/database"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Inventory Management Mini System API
// @version 1.0
// @description Inventory Management Mini System API
// @host localhost:8080
// @BasePath /

func main() {

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize database
	db, err := database.InitPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.Close()

	// Initialize Redis
	redisClient, err := cache.InitRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Error initializing Redis: %v", err)
	}
	defer redisClient.Close()

	// Run migrations
	database.RunMigrations(db)

	// Setup router
	router := gin.Default()

	// Setup middleware
	// authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)
	// cacheMiddleware := middleware.NewCacheMiddleware(redisClient)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	// productHandler := handlers.NewProductHandler(db, redisClient)
	// stockHandler := handlers.NewStockHandler(db, redisClient)

	// Public routers
	public := router.Group("/api/v1")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.GET("/auth/google", authHandler.GoogleLogin)
		public.GET("/auth/google/callback", authHandler.GoogleCallback)
		public.GET("/health", handlers.HealthCheck)
	}

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	fmt.Println("Server started on port", cfg.ServerPort)
}
