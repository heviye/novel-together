package service

import (
	"github.com/heviye/novel-together-backend/internal/middleware"
	"github.com/heviye/novel-together-backend/internal/models"
	"gorm.io/gorm"
	"fmt"
	"os"
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

func (s *NovelService) GetAllNovelsWithStats() ([]map[string]interface{}, error) {
	// disabled panic for debugging
	// 同时打印到 stdout
	fmt.Fprintf(os.Stdout, "[TRACE] GetAllNovelsWithStats called\n")
	rows, err := s.db.Raw(`		SELECT 
			n.id,
			n.title,
			n.is_mainline,
			(SELECT COUNT(*) FROM chapters c WHERE c.novel_id = n.id) as chapter_count,
			(SELECT COUNT(*) FROM likes l JOIN chapters c2 ON l.chapter_id = c2.id WHERE c2.novel_id = n.id) as total_likes,
			(SELECT COUNT(*) FROM comments cm JOIN chapters c3 ON cm.chapter_id = c3.id WHERE c3.novel_id = n.id) as comment_count
		FROM novels n
		ORDER BY total_likes DESC
	`).Rows()
	if err != nil {
		os.WriteFile("/tmp/hermes_debug.log", []byte(fmt.Sprintf("Raw error: %v\n", err)), 0644)
		return nil, err
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	os.WriteFile("/tmp/hermes_debug.log", []byte(fmt.Sprintf("Columns: %v\n", cols)), 0644)
	fmt.Fprintf(os.Stdout, "[TRACE] Columns: %v\n", cols)
	var results []map[string]interface{}
	for rows.Next() {
		var (
			id           string
			title        string
			isMainline   bool
			chapterCount int64
			totalLikes   int64
			commentCount int64
		)
		if err := rows.Scan(&id, &title, &isMainline, &chapterCount, &totalLikes, &commentCount); err != nil {
			fmt.Printf("[DEBUG] Scan error: %v\n", err)
			continue
		}
		os.WriteFile("/tmp/hermes_debug.log", []byte(fmt.Sprintf("row: isMainline=%v chapterCount=%d\n", isMainline, chapterCount)), 0644)
	fmt.Fprintf(os.Stdout, "[TRACE] row: isMainline=%v chapterCount=%d\n", isMainline, chapterCount)
		results = append(results, map[string]interface{}{
			"id":            id,
			"title":         title,
			"is_mainline":   isMainline,
			"chapter_count": chapterCount,
			"total_likes":   totalLikes,
			"comment_count": commentCount,
		})
	}
	if err := rows.Err(); err != nil {
		os.WriteFile("/tmp/hermes_debug.log", []byte(fmt.Sprintf("rows.Err: %v\n", err)), 0644)
		return nil, err
	}
	os.WriteFile("/tmp/hermes_debug.log", []byte(fmt.Sprintf("rows.Next done count=%d\n", len(results))), 0644)
	return results, nil
}

func (s *NovelService) GetNovelWithStats(novelID string) (*models.NovelWithStats, error) {
	novel, err := s.GetByID(novelID)
	if err != nil {
		os.WriteFile("/tmp/hermes_debug.log", []byte(fmt.Sprintf("Raw error: %v\n", err)), 0644)
		return nil, err
	}

	var likeCount int64
	s.db.Model(&models.Like{}).
		Joins("JOIN chapters ON chapters.id = likes.chapter_id").
		Where("chapters.novel_id = ?", novelID).
		Count(&likeCount)

	return &models.NovelWithStats{
		Novel:      *novel,
		TotalLikes: likeCount,
		LikeCount:  likeCount,
	}, nil
}

func (s *NovelService) RecalculateMainline() error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 重置所有
		tx.Model(&models.Novel{}).Update("is_mainline", false)

		// 找出点赞最多的小说
		var maxNovelID string
		var maxLikes int64

		rows, err := tx.Raw(`
			SELECT chapters.novel_id, COUNT(likes.id) as like_count
			FROM likes
			JOIN chapters ON chapters.id = likes.chapter_id
			GROUP BY chapters.novel_id
			ORDER BY like_count DESC
			LIMIT 1
		`).Rows()
		if err != nil {
			return err
		}
		defer rows.Close()

		if rows.Next() {
			var novelID string
			var likeCount int64
			if err := rows.Scan(&novelID, &likeCount); err == nil {
				maxNovelID = novelID
				maxLikes = likeCount
			}
		}

		if maxNovelID != "" {
			tx.Model(&models.Novel{}).Where("id = ?", maxNovelID).Update("is_mainline", true)
		}

		fmt.Printf("[DEBUG] MaxLikes was: %d\n", maxLikes)
		return nil
	})
}
