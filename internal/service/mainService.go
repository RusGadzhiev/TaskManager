package service

import "time"

const (
	Description        = "description"
	Executor           = "executor"
	UserName           = "username"
	TaskId             = "taskId"
	Password           = "password"
	CookieName         = "session_id"
	FilterAllTasks     = "AllTasks"
	FilterMyTasks      = "MyTasks"
	FilterCreatedTasks = "CreatedTasks"
	FilterAssign       = "Assign"
	FilterUnassign     = "Unassign"
	FilterComplete     = "Complete"
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

type service struct {
	UsersService
	SessionsService
	TasksService
}

func NewService(usersService UsersService, sessionsService SessionsService, tasksService TasksService) *service {
	return &service{
		usersService,
		sessionsService,
		tasksService,
	}
}
