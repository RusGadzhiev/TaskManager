## TaskManager

<br/>

## Технологии

- Авторизация с использованием Cookie и Middleware pattern
- API в соответствии с принципами REST
- Хранение Cookie с помощью Redis
- Хранение пользовательских данных с помощью MongoDB
- MySQL для хранения задач пользователей
- Развертывание приложения независимо от платформы с помощью Docker

Технологии: Go, SQL, NoSQL, Git, Docker, Linux, MySQL, MongoDB, Redis, MVC, SOLID, REST, Zap

## Рабочее дерево
```
TaskManager
├── cmd
│   └── task_manager
│       └── main.go
├── config.yaml
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── service
│   │   ├── mainService.go
│   │   ├── sessionsService.go
│   │   ├── tasksService.go
│   │   └── usersService.go
│   ├── storage
│   │   ├── sessionsStorage
│   │   │   └── redis
│   │   │       └── repo.go
│   │   ├── tasksStorage
│   │   │   └── mysql
│   │   │       └── repo.go
│   │   └── usersStorage
│   │       └── mongo
│   │           └── repo.go
│   └── transport
│       └── http
│           ├── httpHandler
│           │   ├── middleware.go
│           │   └── htppHandler.go
│           └── httpServer
│               └── httpServer.go
├── Makefile
├── pkg
│   └── logger
│       └── logger.go
├── templates
│   ├── assign.html
│   ├── complete.html
│   ├── create.html
│   ├── login.html
│   ├── logout.html
│   ├── registration.html
│   └── unassign.html
├── README.md
├── TASK.md
└── Dockerfile

```

##  Начало работы

0. Установите Go, Docker и прочее

1. Клонируйте репозиторий

```bash
git clone https://github.com/RusGadzhiev/TaskManager
```   

2. Запустите контейнеры, используя следующие команды:
```
 make
```

## Контакты 

Гаджиев Руслан - [@RusGadzhiev](https://t.me/driveinyourheart)

