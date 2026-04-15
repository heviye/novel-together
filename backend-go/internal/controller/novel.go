package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heviye/novel-together-backend/internal/service"
)

type NovelController struct {
	novelSvc *service.NovelService
}

func NewNovelController(novelSvc *service.NovelService) *NovelController {
	return &NovelController{novelSvc: novelSvc}
}

func (c *NovelController) List(ctx *gin.Context) {
	page := parseInt(ctx.DefaultQuery("page", "1"))
	limit := parseInt(ctx.DefaultQuery("limit", "20"))

	novels, err := c.novelSvc.List(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load novels"})
		return
	}
	ctx.JSON(http.StatusOK, novels)
}

type CreateNovelInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (c *NovelController) Create(ctx *gin.Context) {
	var input CreateNovelInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorID := ctx.GetString("user_id")
	novel, err := c.novelSvc.Create(service.CreateNovelInput{
		Title:       input.Title,
		Description: input.Description,
		AuthorID:    authorID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create novel"})
		return
	}
	ctx.JSON(http.StatusOK, novel)
}

func (c *NovelController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	novel, err := c.novelSvc.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Novel not found"})
		return
	}
	ctx.JSON(http.StatusOK, novel)
}

func (c *NovelController) GetChapters(ctx *gin.Context) {
	id := ctx.Param("id")
	chapters, err := c.novelSvc.GetChapters(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get chapters"})
		return
	}
	ctx.JSON(http.StatusOK, chapters)
}
