package service

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
