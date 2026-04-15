package main

import (
	"flag"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/heviye/novel-together-backend/internal/config"
	"github.com/heviye/novel-together-backend/internal/middleware"
	"github.com/heviye/novel-together-backend/internal/models"
	"github.com/heviye/novel-together-backend/internal/routes"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	// Load config
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Check JWT_SECRET
	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required in config")
	}

	// Set JWT secret for middleware
	middleware.SetJWTSecret(cfg.JWT.Secret)

	// Initialize database
	db, err := models.InitDBWithDSN(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(
		&models.User{},
		&models.Novel{},
		&models.Chapter{},
		&models.Like{},
		&models.Comment{},
		&models.Follow{},
	); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// Setup Gin
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.Abort()
			return
		}
		c.Next()
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API routes
	api := r.Group("/api")
	{
		routes.RegisterAuthRoutes(api, db)
		routes.RegisterUserRoutes(api, db)
		routes.RegisterNovelRoutes(api, db)
		routes.RegisterChapterRoutes(api, db)
	}

	port := cfg.App.Port
	if port == "" {
		port = "3000"
	}
	host := cfg.App.Host
	if host == "" {
		host = "0.0.0.0"
	}
	r.Run(host + ":" + port)
}
