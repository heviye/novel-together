package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/heviye/novel-together-backend/internal/controller"
	"github.com/heviye/novel-together-backend/internal/middleware"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(api *gin.RouterGroup, db *gorm.DB, authCtrl *controller.AuthController) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", authCtrl.Register)
		auth.POST("/login", authCtrl.Login)
	}
}

func RegisterUserRoutes(api *gin.RouterGroup, db *gorm.DB, userCtrl *controller.UserController) {
	users := api.Group("/users")
	{
		users.GET("/:id", userCtrl.GetProfile)

		secured := users.Group("")
		secured.Use(middleware.AuthMiddleware())
		{
			secured.PUT("/:id", userCtrl.UpdateProfile)
			secured.POST("/:id/follow", userCtrl.Follow)
			secured.DELETE("/:id/follow", userCtrl.Unfollow)
		}

		users.GET("/:id/followers", userCtrl.GetFollowers)
		users.GET("/:id/following", userCtrl.GetFollowing)
	}
}

func RegisterNovelRoutes(api *gin.RouterGroup, db *gorm.DB, novelCtrl *controller.NovelController) {
	novels := api.Group("/novels")
	{
		novels.GET("", novelCtrl.List)

		secured := novels.Group("")
		secured.Use(middleware.AuthMiddleware())
		{
			secured.POST("", novelCtrl.Create)
		}

		novels.GET("/:id", novelCtrl.Get)
		novels.GET("/:id/chapters", novelCtrl.GetChapters)
		novels.GET("/:id/stats", novelCtrl.GetStats)
	}

	// 管理接口（临时开放，后续可加权限）
	admin := api.Group("/admin/novels")
	{
		admin.POST("/recalculate-mainline", novelCtrl.RecalculateMainline)
		admin.GET("/all-stats", novelCtrl.GetAllStats)
	}
}

func RegisterChapterRoutes(api *gin.RouterGroup, db *gorm.DB, chapterCtrl *controller.ChapterController) {
	chapters := api.Group("/chapters")

	secured := chapters.Group("")
	secured.Use(middleware.AuthMiddleware())
	{
		secured.POST("/novels/:novelId/chapters", chapterCtrl.Create)
		secured.POST("/:id/like", chapterCtrl.Like)
		secured.DELETE("/:id/like", chapterCtrl.Unlike)
		secured.POST("/:id/comments", chapterCtrl.Comment)
	}

	chapters.GET("/:id", chapterCtrl.Get)
	chapters.GET("/:id/likes", chapterCtrl.GetLikes)
	chapters.GET("/:id/comments", chapterCtrl.GetComments)
}
