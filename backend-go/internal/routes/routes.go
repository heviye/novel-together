package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heviye/novel-together-backend/internal/middleware"
	"github.com/heviye/novel-together-backend/internal/models"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(api *gin.RouterGroup, db *gorm.DB) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			var input struct {
				Username string `json:"username" binding:"required"`
				Email    string `json:"email" binding:"required,email"`
				Password string `json:"password" binding:"required,min=6"`
			}
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			var existing models.User
			if err := db.Where("email = ? OR username = ?", input.Email, input.Username).First(&existing).Error; err == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Username or email already exists"})
				return
			}

			hashed, _ := middleware.HashPassword(input.Password)
			user := models.User{
				ID:       middleware.GenerateUUID(),
				Username: input.Username,
				Email:    input.Email,
				Password: hashed,
			}
			if err := db.Create(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"id": user.ID, "username": user.Username, "email": user.Email})
		})

		auth.POST("/login", func(c *gin.Context) {
			var input struct {
				Email    string `json:"email" binding:"required,email"`
				Password string `json:"password" binding:"required"`
			}
			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			var user models.User
			if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			if !middleware.CheckPasswordHash(input.Password, user.Password) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				return
			}

			token, _ := middleware.GenerateToken(user.ID, user.Username)
			c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "username": user.Username, "email": user.Email}})
		})
	}
}

func RegisterUserRoutes(api *gin.RouterGroup, db *gorm.DB) {
	users := api.Group("/users")
	{
		users.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			var user models.User
			if err := db.First(&user, "id = ?", id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"id": user.ID, "username": user.Username, "bio": user.Bio, "avatar_url": user.AvatarURL, "created_at": user.CreatedAt})
		})

		secured := users.Group("")
		secured.Use(middleware.AuthMiddleware())
		{
			secured.PUT("/:id", func(c *gin.Context) {
				userID := c.GetString("user_id")
				id := c.Param("id")
				if userID != id {
					c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
					return
				}

				var input struct {
					Bio       *string `json:"bio"`
					AvatarURL *string `json:"avatar_url"`
				}
				c.ShouldBindJSON(&input)

				db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{"bio": input.Bio, "avatar_url": input.AvatarURL})
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			secured.POST("/:id/follow", func(c *gin.Context) {
				userID := c.GetString("user_id")
				id := c.Param("id")
				if userID == id {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot follow yourself"})
					return
				}

				follow := models.Follow{
					ID:          middleware.GenerateUUID(),
					FollowerID:  userID,
					FollowingID: id,
				}
				db.Where("follower_id = ? AND following_id = ?", userID, id).Delete(&models.Follow{})
				db.Create(&follow)
				c.JSON(http.StatusOK, gin.H{"success": true})
			})

			secured.DELETE("/:id/follow", func(c *gin.Context) {
				userID := c.GetString("user_id")
				id := c.Param("id")
				db.Where("follower_id = ? AND following_id = ?", userID, id).Delete(&models.Follow{})
				c.JSON(http.StatusOK, gin.H{"success": true})
			})
		}

		users.GET("/:id/followers", func(c *gin.Context) {
			id := c.Param("id")
			var followers []models.User
			db.Table("users").Select("users.id, users.username, users.avatar_url").
				Joins("JOIN follows ON follows.follower_id = users.id").
				Where("follows.following_id = ?", id).
				Scan(&followers)
			c.JSON(http.StatusOK, followers)
		})

		users.GET("/:id/following", func(c *gin.Context) {
			id := c.Param("id")
			var following []models.User
			db.Table("users").Select("users.id, users.username, users.avatar_url").
				Joins("JOIN follows ON follows.following_id = users.id").
				Where("follows.follower_id = ?", id).
				Scan(&following)
			c.JSON(http.StatusOK, following)
		})
	}
}

func RegisterNovelRoutes(api *gin.RouterGroup, db *gorm.DB) {
	novels := api.Group("/novels")
	{
		novels.GET("", func(c *gin.Context) {
			page := c.DefaultQuery("page", "1")
			limit := c.DefaultQuery("limit", "20")

			var novels []models.Novel
			db.Preload("Author", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, username, avatar_url")
			}).
				Order("updated_at DESC").
				Offset((parseInt(page) - 1) * parseInt(limit)).
				Limit(parseInt(limit)).
				Find(&novels)

			c.JSON(http.StatusOK, novels)
		})

		secured := novels.Group("")
		secured.Use(middleware.AuthMiddleware())
		{
			secured.POST("", func(c *gin.Context) {
				userID := c.GetString("user_id")
				var input struct {
					Title       string `json:"title" binding:"required"`
					Description string `json:"description"`
				}
				c.ShouldBindJSON(&input)

				novel := models.Novel{
					ID:          middleware.GenerateUUID(),
					Title:       input.Title,
					Description: strPtr(input.Description),
					AuthorID:    userID,
					Status:      "active",
				}
				db.Create(&novel)
				c.JSON(http.StatusOK, novel)
			})
		}

		novels.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")
			var novel models.Novel
			if err := db.Preload("Author", "id, username").First(&novel, "id = ?", id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Novel not found"})
				return
			}
			c.JSON(http.StatusOK, novel)
		})

		novels.GET("/:id/chapters", func(c *gin.Context) {
			id := c.Param("id")
			var chapters []models.Chapter
			db.Preload("Author", "id, username").
				Where("novel_id = ?", id).
				Order("chapter_number").
				Find(&chapters)
			c.JSON(http.StatusOK, chapters)
		})
	}
}

func RegisterChapterRoutes(api *gin.RouterGroup, db *gorm.DB) {
	chapters := api.Group("/chapters")

	secured := chapters.Group("")
	secured.Use(middleware.AuthMiddleware())
	{
		secured.POST("/novels/:novelId/chapters", func(c *gin.Context) {
			userID := c.GetString("user_id")
			novelID := c.Param("novelId")

			var novel models.Novel
			if err := db.First(&novel, "id = ?", novelID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Novel not found"})
				return
			}

			var input struct {
				Content string `json:"content" binding:"required"`
			}
			c.ShouldBindJSON(&input)

			var maxChapter models.Chapter
			db.Model(&models.Chapter{}).Where("novel_id = ?", novelID).Select("MAX(chapter_number)").Scan(&maxChapter)

			chapter := models.Chapter{
				ID:            middleware.GenerateUUID(),
				NovelID:       novelID,
				ChapterNumber: maxChapter.ChapterNumber + 1,
				AuthorID:      userID,
				Content:       input.Content,
			}
			db.Create(&chapter)

			db.Model(&models.Novel{}).Where("id = ?", novelID).Update("updated_at", middleware.Now())

			c.JSON(http.StatusOK, chapter)
		})

		secured.POST("/:id/like", func(c *gin.Context) {
			userID := c.GetString("user_id")
			chapterID := c.Param("id")

			like := models.Like{
				ID:        middleware.GenerateUUID(),
				UserID:    userID,
				ChapterID: chapterID,
			}
			db.Where("user_id = ? AND chapter_id = ?", userID, chapterID).Delete(&models.Like{})
			db.Create(&like)
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		secured.DELETE("/:id/like", func(c *gin.Context) {
			userID := c.GetString("user_id")
			chapterID := c.Param("id")
			db.Where("user_id = ? AND chapter_id = ?", userID, chapterID).Delete(&models.Like{})
			c.JSON(http.StatusOK, gin.H{"success": true})
		})

		secured.POST("/:id/comments", func(c *gin.Context) {
			userID := c.GetString("user_id")
			chapterID := c.Param("id")

			var input struct {
				Content string `json:"content" binding:"required"`
			}
			c.ShouldBindJSON(&input)

			comment := models.Comment{
				ID:        middleware.GenerateUUID(),
				UserID:    userID,
				ChapterID: chapterID,
				Content:   input.Content,
			}
			db.Create(&comment)
			c.JSON(http.StatusOK, comment)
		})
	}

	chapters.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		var chapter models.Chapter
		if err := db.Preload("Author", "id, username").Preload("Novel", "id, title").First(&chapter, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
			return
		}
		c.JSON(http.StatusOK, chapter)
	})

	chapters.GET("/:id/likes", func(c *gin.Context) {
		id := c.Param("id")
		var count int64
		db.Model(&models.Like{}).Where("chapter_id = ?", id).Count(&count)
		c.JSON(http.StatusOK, gin.H{"count": count})
	})

	chapters.GET("/:id/comments", func(c *gin.Context) {
		id := c.Param("id")
		var comments []models.Comment
		db.Preload("User", "id, username").
			Where("chapter_id = ?", id).
			Order("created_at").
			Find(&comments)
		c.JSON(http.StatusOK, comments)
	})
}

func parseInt(s string) int {
	var n int
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}

func strPtr(s string) *string {
	return &s
}
