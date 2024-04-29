package service

import (
	"time"
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
	Dur       time.Duration
}
