package service

import (
	"github.com/heviye/novel-together-backend/internal/middleware"
	"github.com/heviye/novel-together-backend/internal/models"
	"gorm.io/gorm"
)

type ChapterService struct {
	db *gorm.DB
}

func NewChapterService(db *gorm.DB) *ChapterService {
	return &ChapterService{db: db}
}

type CreateChapterInput struct {
	NovelID  string `json:"novel_id"`
	AuthorID string `json:"author_id"`
	Content  string `json:"content"`
}

func (s *ChapterService) Create(input CreateChapterInput) (*models.Chapter, error) {
	var novel models.Novel
	if err := s.db.First(&novel, "id = ?", input.NovelID).Error; err != nil {
		return nil, ErrNovelNotFound
	}

	// Get max chapter number
	var maxChapter models.Chapter
	s.db.Model(&models.Chapter{}).Where("novel_id = ?", input.NovelID).Select("MAX(chapter_number)").Scan(&maxChapter)

	chapter := models.Chapter{
		ID:            middleware.GenerateUUID(),
		NovelID:       input.NovelID,
		ChapterNumber: maxChapter.ChapterNumber + 1,
		AuthorID:      input.AuthorID,
		Content:       input.Content,
	}
	if err := s.db.Create(&chapter).Error; err != nil {
		return nil, err
	}

	// Update novel timestamp
	s.db.Model(&models.Novel{}).Where("id = ?", input.NovelID).Update("updated_at", middleware.Now())

	return &chapter, nil
}

func (s *ChapterService) GetByID(id string) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := s.db.Preload("Author", "id, username").Preload("Novel", "id, title").First(&chapter, "id = ?", id).Error; err != nil {
		return nil, ErrChapterNotFound
	}
	return &chapter, nil
}

type LikeInput struct {
	UserID    string `json:"user_id"`
	ChapterID string `json:"chapter_id"`
}

func (s *ChapterService) Like(input LikeInput) error {
	like := models.Like{
		ID:        middleware.GenerateUUID(),
		UserID:    input.UserID,
		ChapterID: input.ChapterID,
	}
	s.db.Where("user_id = ? AND chapter_id = ?", input.UserID, input.ChapterID).Delete(&models.Like{})
	return s.db.Create(&like).Error
}

func (s *ChapterService) Unlike(input LikeInput) error {
	return s.db.Where("user_id = ? AND chapter_id = ?", input.UserID, input.ChapterID).Delete(&models.Like{}).Error
}

func (s *ChapterService) GetLikeCount(chapterID string) (int64, error) {
	var count int64
	err := s.db.Model(&models.Like{}).Where("chapter_id = ?", chapterID).Count(&count).Error
	return count, err
}

type CommentInput struct {
	UserID    string `json:"user_id"`
	ChapterID string `json:"chapter_id"`
	Content   string `json:"content"`
}

func (s *ChapterService) Comment(input CommentInput) (*models.Comment, error) {
	comment := models.Comment{
		ID:        middleware.GenerateUUID(),
		UserID:    input.UserID,
		ChapterID: input.ChapterID,
		Content:   input.Content,
	}
	if err := s.db.Create(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (s *ChapterService) GetComments(chapterID string) ([]models.Comment, error) {
	var comments []models.Comment
	err := s.db.Preload("User", "id, username").
		Where("chapter_id = ?", chapterID).
		Order("created_at").
		Find(&comments).Error
	return comments, err
}
