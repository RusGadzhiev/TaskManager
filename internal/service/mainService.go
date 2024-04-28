package service


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
