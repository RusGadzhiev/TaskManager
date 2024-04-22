package sessions

import (
	"context"
	"time"
)

// хранит мапу userId -> cookie
type SessionsRepo interface {
	// возвращает username пользователя по значению сессии
	GetUser(ctx context.Context, cookieVal string) (string, error)
	// добавляет новую сессию
	Add(ctx context.Context, cookieVal string, username string, dur time.Duration) error
	// удаляет cookie
	Delete(ctx context.Context, cookieVal string) error
}
