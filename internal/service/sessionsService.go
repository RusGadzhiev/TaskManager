package service

import (
	"context"
	"math/rand"
	"time"
)

type SessionsStorage interface {
	// возвращает username пользователя по значению сессии
	GetUser(ctx context.Context, cookieVal string) (string, error)
	// добавляет новую сессию
	Add(ctx context.Context, cookieVal string, username string, dur time.Duration) error
	// удаляет cookie
	Delete(ctx context.Context, cookieVal string) error
}

type SessionsService struct {
	repo SessionsStorage
}

func NewSessionsService(repo SessionsStorage) *SessionsService {
	return &SessionsService{
		repo: repo,
	}
}

func (s *SessionsService) DeleteCookie(ctx context.Context, cookieVal string) error {
	return s.repo.Delete(ctx, cookieVal)
}

func (s *SessionsService) AddCookie(ctx context.Context, username string) (*Session, error) {
	cookieVal := randStringChars(32)
	dur :=  72 * time.Hour
	err := s.repo.Add(ctx, cookieVal, username, dur)
	if err != nil {
		return nil, err
	}
	session := &Session{
		CookieVal: cookieVal,
		Dur: dur,
	}
	return session, nil
}

func (s *SessionsService) GetUserByCookie(ctx context.Context, cookieVal string) (string, error) {
	return s.repo.GetUser(ctx, cookieVal)
}

func randStringChars(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
