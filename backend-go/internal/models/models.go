package models

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Bio       *string   `json:"bio"`
	AvatarURL *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Novel struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description *string   `json:"description"`
	AuthorID    string    `gorm:"not null" json:"author_id"`
	Author      User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	IsMainline  bool      `gorm:"default:false" json:"is_mainline"`
	Status      string    `gorm:"default:active" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NovelWithStats struct {
	Novel
	TotalLikes int64 `json:"total_likes"`
	LikeCount  int64 `json:"like_count"`
}


type Chapter struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	NovelID       string    `gorm:"not null" json:"novel_id"`
	Novel         Novel     `gorm:"foreignKey:NovelID" json:"novel,omitempty"`
	ChapterNumber int       `gorm:"not null" json:"chapter_number"`
	AuthorID      string    `gorm:"not null" json:"author_id"`
	Author        User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Content       string    `gorm:"not null" json:"content"`
	CreatedAt     time.Time `json:"created_at"`
}

type Like struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"not null" json:"user_id"`
	ChapterID string    `gorm:"not null" json:"chapter_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ChapterID string    `gorm:"not null" json:"chapter_id"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Follow struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	FollowerID  string    `gorm:"not null" json:"follower_id"`
	FollowingID string    `gorm:"not null" json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func InitDBWithDSN(dsn string, driver string) (*gorm.DB, error) {
	var dialect gorm.Dialector
	switch driver {
	case "sqlite":
		dialect = sqlite.Open(dsn)
	case "postgres":
		dialect = postgres.Open(dsn)
	default:
		dialect = sqlite.Open(dsn)
	}

	db, err := gorm.Open(dialect, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
