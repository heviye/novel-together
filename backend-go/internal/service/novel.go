package service

import (
	"github.com/heviye/novel-together-backend/internal/middleware"
	"github.com/heviye/novel-together-backend/internal/models"
	"gorm.io/gorm"
)

type NovelService struct {
	db *gorm.DB
}

func NewNovelService(db *gorm.DB) *NovelService {
	return &NovelService{db: db}
}

type CreateNovelInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorID    string `json:"author_id"`
}

func (s *NovelService) Create(input CreateNovelInput) (*models.Novel, error) {
	novel := models.Novel{
		ID:          middleware.GenerateUUID(),
		Title:       input.Title,
		Description: StrPtr(input.Description),
		AuthorID:    input.AuthorID,
		Status:      "active",
	}
	if err := s.db.Create(&novel).Error; err != nil {
		return nil, err
	}
	return &novel, nil
}

func (s *NovelService) GetByID(id string) (*models.Novel, error) {
	var novel models.Novel
	if err := s.db.Preload("Author", "id, username").First(&novel, "id = ?", id).Error; err != nil {
		return nil, ErrNovelNotFound
	}
	return &novel, nil
}

func (s *NovelService) List(page, limit int) ([]models.Novel, error) {
	var novels []models.Novel
	offset := (page - 1) * limit
	err := s.db.Preload("Author", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username, avatar_url")
	}).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&novels).Error
	return novels, err
}

func (s *NovelService) GetChapters(novelID string) ([]models.Chapter, error) {
	var chapters []models.Chapter
	err := s.db.Preload("Author", "id, username").
		Where("novel_id = ?", novelID).
		Order("chapter_number").
		Find(&chapters).Error
	return chapters, err
}
