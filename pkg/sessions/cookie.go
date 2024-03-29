package sessions

import (
	"HW4/pkg/users"
	"context"
	"net/http"
)

type Session struct {
	Id     uint64      `json:"id"`
	Cookie http.Cookie `json:"cookie"`
	User   users.User  `json:"user"`
}

type SessionsRepo interface {
	GetSession(ctx context.Context, id uint64) (*Session, error)
	Add(ctx context.Context, session *Session) error
}
