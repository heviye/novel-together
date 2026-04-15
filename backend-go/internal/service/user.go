package service

import (
	"github.com/heviye/novel-together-backend/internal/middleware"
	"github.com/heviye/novel-together-backend/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *UserService) Register(input RegisterInput) (*models.User, error) {
	var existing models.User
	if err := s.db.Where("email = ? OR username = ?", input.Email, input.Username).First(&existing).Error; err == nil {
		return nil, ErrUserExists
	}

	hashed, err := middleware.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		ID:       middleware.GenerateUUID(),
		Username: input.Username,
		Email:    input.Email,
		Password: hashed,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOutput struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

func (s *UserService) Login(input LoginInput) (*LoginOutput, error) {
	var user models.User
	if err := s.db.Where("email = ?", input.Email).First(&user).Error; err != nil {
		return nil, ErrInvalidCredentials
	}

	if !middleware.CheckPasswordHash(input.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token: token,
		User:  &user,
	}, nil
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

type UpdateUserInput struct {
	Bio       *string `json:"bio"`
	AvatarURL *string `json:"avatar_url"`
}

func (s *UserService) Update(id string, input UpdateUserInput) error {
	result := s.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"bio":        input.Bio,
		"avatar_url": input.AvatarURL,
	})
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *UserService) Follow(followerID, followingID string) error {
	if followerID == followingID {
		return ErrCannotFollowSelf
	}

	follow := models.Follow{
		ID:          middleware.GenerateUUID(),
		FollowerID:  followerID,
		FollowingID: followingID,
	}
	s.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&models.Follow{})
	return s.db.Create(&follow).Error
}

func (s *UserService) Unfollow(followerID, followingID string) error {
	return s.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&models.Follow{}).Error
}

func (s *UserService) GetFollowers(id string) ([]models.User, error) {
	var users []models.User
	err := s.db.Table("users").Select("users.id, users.username, users.avatar_url").
		Joins("JOIN follows ON follows.follower_id = users.id").
		Where("follows.following_id = ?", id).
		Scan(&users).Error
	return users, err
}

func (s *UserService) GetFollowing(id string) ([]models.User, error) {
	var users []models.User
	err := s.db.Table("users").Select("users.id, users.username, users.avatar_url").
		Joins("JOIN follows ON follows.following_id = users.id").
		Where("follows.follower_id = ?", id).
		Scan(&users).Error
	return users, err
}
