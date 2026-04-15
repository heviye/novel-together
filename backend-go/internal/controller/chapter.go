package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/heviye/novel-together-backend/internal/service"
)

type ChapterController struct {
	chapterSvc *service.ChapterService
}

func NewChapterController(chapterSvc *service.ChapterService) *ChapterController {
	return &ChapterController{chapterSvc: chapterSvc}
}

type CreateChapterInput struct {
	Content string `json:"content"`
}

func (c *ChapterController) Create(ctx *gin.Context) {
	novelID := ctx.Param("novelId")
	userID := ctx.GetString("user_id")

	var input CreateChapterInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chapter, err := c.chapterSvc.Create(service.CreateChapterInput{
		NovelID:  novelID,
		AuthorID: userID,
		Content:  input.Content,
	})
	if err != nil {
		if err == service.ErrNovelNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chapter"})
		return
	}
	ctx.JSON(http.StatusOK, chapter)
}

func (c *ChapterController) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	chapter, err := c.chapterSvc.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Chapter not found"})
		return
	}
	ctx.JSON(http.StatusOK, chapter)
}

func (c *ChapterController) Like(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")

	if err := c.chapterSvc.Like(service.LikeInput{UserID: userID, ChapterID: id}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Like failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (c *ChapterController) Unlike(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")

	if err := c.chapterSvc.Unlike(service.LikeInput{UserID: userID, ChapterID: id}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Unlike failed"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (c *ChapterController) GetLikes(ctx *gin.Context) {
	id := ctx.Param("id")
	count, err := c.chapterSvc.GetLikeCount(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get likes"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"count": count})
}

type CreateCommentInput struct {
	Content string `json:"content"`
}

func (c *ChapterController) Comment(ctx *gin.Context) {
	id := ctx.Param("id")
	userID := ctx.GetString("user_id")

	var input CreateCommentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := c.chapterSvc.Comment(service.CommentInput{
		UserID:    userID,
		ChapterID: id,
		Content:   input.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to comment"})
		return
	}
	ctx.JSON(http.StatusOK, comment)
}

func (c *ChapterController) GetComments(ctx *gin.Context) {
	id := ctx.Param("id")
	comments, err := c.chapterSvc.GetComments(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}
	ctx.JSON(http.StatusOK, comments)
}
