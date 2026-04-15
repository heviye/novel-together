package service

import "errors"

var (
	ErrUserExists         = errors.New("username or email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrCannotFollowSelf   = errors.New("cannot follow yourself")
	ErrNovelNotFound      = errors.New("novel not found")
	ErrChapterNotFound    = errors.New("chapter not found")
)
