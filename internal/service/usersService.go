package service

import (
	"context"
	"errors"
)

var (
	ErrNoUser            = errors.New("no such user")
	ErrUserExist         = errors.New("user with this login exists")
	ErrIncorrectPassword = errors.New("incorrect password")
)

type UsersStorage interface {
	// возвращает ошибку service.ErrNoUser если юзера нет
	GetUser(ctx context.Context, username string) (*User, error)
	// возвращает ошибку service.ErrUserExist если юзер уже есть
	AddUser(ctx context.Context, user *User) error
}

type UsersService struct {
	repo UsersStorage
}

func NewUsersService(repo UsersStorage) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s *UsersService) Authentificate(ctx context.Context, user *User) error {
	realUser, err := s.repo.GetUser(ctx, user.UserName)
	if err != nil {
		return err
	}

	if user.Password != realUser.Password {
		return ErrIncorrectPassword
	}

	return nil
}

func (s *UsersService) AddUser(ctx context.Context, user *User) error {
	return s.repo.AddUser(ctx, user)
}
