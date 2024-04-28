package service


type service struct {
	usersService UsersService
	sessionsService SessionsService
	tasksService    TasksService
}

func NewService(usersService UsersService, sessionsService SessionsService, tasksService TasksService) *service {
	return &service{
		usersService: usersService,
		sessionsService: sessionsService,
		tasksService: tasksService,
	}
}
