package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heviye/novel-together-backend/internal/service"
)

type UserController struct {
	userSvc *service.UserService
}

func NewUserController(userSvc *service.UserService) *UserController {
	return &UserController{userSvc: userSvc}
}

func (c *UserController) GetProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := c.userSvc.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": user.ID, "username": user.Username, "bio": user.Bio, "avatar_url": user.AvatarURL, "created_at": user.CreatedAt})
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")
	if userID != id {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	var input service.UpdateUserInput
	ctx.ShouldBindJSON(&input)

	if err := c.userSvc.Update(id, input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (c *UserController) Follow(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")

	if err := c.userSvc.Follow(userID, id); err != nil {
		if err == service.ErrCannotFollowSelf {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Follow failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (c *UserController) Unfollow(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")

	if err := c.userSvc.Unfollow(userID, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unfollow failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (c *UserController) GetFollowers(ctx *gin.Context) {
	id := ctx.Param("id")
	users, err := c.userSvc.GetFollowers(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get followers"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *UserController) GetFollowing(ctx *gin.Context) {
	id := ctx.Param("id")
	users, err := c.userSvc.GetFollowing(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get following"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
