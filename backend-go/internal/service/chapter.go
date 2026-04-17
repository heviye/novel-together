package service

import (
	"fmt"
	"sync"

	"github.com/heviye/novel-together-backend/internal/middleware"
	"github.com/heviye/novel-together-backend/internal/models"
	"gorm.io/gorm"
)

type ChapterService struct {
	db *gorm.DB
	mainlineMutex sync.Mutex
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

	var maxChapterNumber int
	if err := s.db.Model(&models.Chapter{}).Where("novel_id = ?", input.NovelID).
		Select("COALESCE(MAX(chapter_number), 0)").Scan(&maxChapterNumber).Error; err != nil {
		return nil, err
	}

	chapter := models.Chapter{
		ID:            middleware.GenerateUUID(),
		NovelID:       input.NovelID,
		ChapterNumber: maxChapterNumber + 1,
		AuthorID:      input.AuthorID,
		Content:       input.Content,
	}
	if err := s.db.Create(&chapter).Error; err != nil {
		return nil, err
	}

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
	if err := s.db.Create(&like).Error; err != nil {
		return err
	}

	// 异步触发主线重计算
	go s.triggerRecalculateMainline()
	return nil
}

func (s *ChapterService) Unlike(input LikeInput) error {
	if err := s.db.Where("user_id = ? AND chapter_id = ?", input.UserID, input.ChapterID).Delete(&models.Like{}).Error; err != nil {
		return err
	}
	go s.triggerRecalculateMainline()
	return nil
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

// triggerRecalculateMainline 异步触发主线重计算
func (s *ChapterService) triggerRecalculateMainline() {
	s.mainlineMutex.Lock()
	defer s.mainlineMutex.Unlock()

	// 直接在db上执行SQL重计算，避免循环import
	// 找出点赞最多的小说并设置is_mainline=true
	s.db.Transaction(func(tx *gorm.DB) error {
		// 全部设为false
		tx.Model(&models.Novel{}).Update("is_mainline", false)

		// 找出点赞最多的小说
		var maxNovelID string
		var maxLikes int64 = 0

		rows, err := tx.Raw(`SELECT chapters.novel_id, COUNT(likes.id) as cnt
			FROM likes JOIN chapters ON chapters.id = likes.chapter_id
			GROUP BY chapters.novel_id ORDER BY cnt DESC LIMIT 1`).Rows()
		if err != nil {
			fmt.Printf("[ERROR] recalc failed: %v\n", err)
			return err
		}
		defer rows.Close()

		if rows.Next() {
			var novelID string
			var cnt int64
			rows.Scan(&novelID, &cnt)
			maxNovelID = novelID
			maxLikes = cnt
		}

		if maxNovelID != "" {
			tx.Model(&models.Novel{}).Where("id = ?", maxNovelID).Update("is_mainline", true)
			fmt.Printf("[Mainline] Novel %s set as mainline with %d likes\n", maxNovelID, maxLikes)
		}
		return nil
	})
}
