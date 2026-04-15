package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heviye/novel-together-backend/internal/service"
)

type AuthController struct {
	userSvc *service.UserService
}

func NewAuthController(userSvc *service.UserService) *AuthController {
	return &AuthController{userSvc: userSvc}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var input service.RegisterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userSvc.Register(input)
	if err != nil {
		if err == service.ErrUserExists {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": user.ID, "username": user.Username, "email": user.Email})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input service.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := c.userSvc.Login(input)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": output.Token, "user": gin.H{"id": output.User.ID, "username": output.User.Username, "email": output.User.Email}})
}
