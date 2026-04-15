package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/heviye/novel-together-backend/internal/models"
	"github.com/heviye/novel-together-backend/internal/routes"
)

func main() {
	// Check JWT_SECRET
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET is required")
	}

	// Initialize database
	db, err := models.InitDB()
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	r.Run(":" + port)
}
