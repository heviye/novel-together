package models

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Bio       *string   `json:"bio"`
	AvatarURL *string   `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Novel struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Description *string   `json:"description"`
	AuthorID    string    `gorm:"type:uuid;not null" json:"author_id"`
	Author      User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Status      string    `gorm:"default:active" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Chapter struct {
	ID            string    `gorm:"type:uuid;primaryKey" json:"id"`
	NovelID       string    `gorm:"type:uuid;not null" json:"novel_id"`
	Novel         Novel     `gorm:"foreignKey:NovelID" json:"novel,omitempty"`
	ChapterNumber int       `gorm:"not null" json:"chapter_number"`
	AuthorID      string    `gorm:"type:uuid;not null" json:"author_id"`
	Author        User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Content       string    `gorm:"type:text;not null" json:"content"`
	CreatedAt     time.Time `json:"created_at"`
}

type Like struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	ChapterID string    `gorm:"type:uuid;not null" json:"chapter_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ChapterID string    `gorm:"type:uuid;not null" json:"chapter_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Follow struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	FollowerID  string    `gorm:"type:uuid;not null" json:"follower_id"`
	FollowingID string    `gorm:"type:uuid;not null" json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func InitDB() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "novel_together"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
