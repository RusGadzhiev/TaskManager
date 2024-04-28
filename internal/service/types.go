package service

import (
	"errors"
	"time"
)

const (
	FilterAllTasks     = "AllTasks"
	FilterMyTasks      = "MyTasks"
	FilterCreatedTasks = "CreatedTasks"
	FilterAssign       = "Assign"
	FilterUnassign     = "Unassign"
	FilterComplete     = "Complete"
)

const (
	Description = "description"
	Executor    = "executor"
	UserName    = "username"
	TaskId      = "taskId"
	Password    = "password"
	CookieName  = "session_id"
)

var (
	ErrBadId             = errors.New("bad id")
	ErrNoUser            = errors.New("no such user")
	ErrUserExist         = errors.New("user with this login exists")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrNoUserBySession   = errors.New("no user by session")
	letterRunes          = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

type Task struct {
	ID          uint64
	Owner       string
	Executor    string
	Description string
	Completed   bool
	Assigned    bool
}

type User struct {
	UserName string `bson:"username"`
	Password string `bson:"password"`
}

type Session struct {
	CookieVal string
	Dur time.Duration
}
