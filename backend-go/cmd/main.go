package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/heviye/novel-together-backend/internal/config"
	"github.com/heviye/novel-together-backend/internal/controller"
	"github.com/heviye/novel-together-backend/internal/middleware"
	"github.com/heviye/novel-together-backend/internal/models"
	"github.com/heviye/novel-together-backend/internal/routes"
	"github.com/heviye/novel-together-backend/internal/service"
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
	db, err := models.InitDBWithDSN(cfg.Database.DSN(), cfg.Database.Driver)
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

	// Create services
	userSvc := service.NewUserService(db)
	novelSvc := service.NewNovelService(db)
	chapterSvc := service.NewChapterService(db)

	// Create controllers
	authCtrl := controller.NewAuthController(userSvc)
	userCtrl := controller.NewUserController(userSvc)
	novelCtrl := controller.NewNovelController(novelSvc)
	chapterCtrl := controller.NewChapterController(chapterSvc)

	// Setup Gin
	r := gin.Default()

	// CORS middleware - configurable via config
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := cfg.CORS.AllowedOrigin(origin)
		if allowed || cfg.CORS.Origins == "" || cfg.CORS.Origins == "*" {
			if origin != "" {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			}
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			method := c.Request.Header.Get("Access-Control-Request-Method")
			if method == "" {
				method = c.Request.Header.Get("Access-Control-Request-Method")
			}
			if strings.EqualFold(method, "POST") || strings.EqualFold(method, "GET") ||
				strings.EqualFold(method, "PUT") || strings.EqualFold(method, "DELETE") ||
				strings.EqualFold(method, "PATCH") || strings.EqualFold(method, "OPTIONS") {
				c.AbortWithStatus(204)
				return
			}
			c.AbortWithStatus(204)
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
		routes.RegisterAuthRoutes(api, db, authCtrl)
		routes.RegisterUserRoutes(api, db, userCtrl)
		routes.RegisterNovelRoutes(api, db, novelCtrl)
		routes.RegisterChapterRoutes(api, db, chapterCtrl)
	}

	port := cfg.App.Port
	if port == "" {
		port = "3000"
	}
	host := cfg.App.Host
	if host == "" {
		host = "0.0.0.0"
	}
	log.Printf("Server starting on %s:%s", host, port)
	r.Run(host + ":" + port)
}
