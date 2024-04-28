package main

import (
	"context"
	"html/template"
	"os/signal"
	"syscall"

	"github.com/RusGadzhiev/TaskManager/internal/config"
	"github.com/RusGadzhiev/TaskManager/internal/service"
	"github.com/RusGadzhiev/TaskManager/internal/storage/sessionsStorage/redis"
	"github.com/RusGadzhiev/TaskManager/internal/storage/tasksStorage/mysql"
	"github.com/RusGadzhiev/TaskManager/internal/storage/usersStorage/mongo"
	"github.com/RusGadzhiev/TaskManager/internal/transport/http/httpHandler"
	"github.com/RusGadzhiev/TaskManager/internal/transport/http/httpServer"
	"go.uber.org/zap"

	"github.com/RusGadzhiev/TaskManager/pkg/logger"
)

// что делать приватным а что публичным
const (
	templatePattern = "./templates/*"
)

type Server interface {
	Run(ctx context.Context, logger *zap.SugaredLogger) error
}

func main() {
	cfg := config.MustLoad()
	tmpl := template.Must(template.ParseGlob(templatePattern))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger := logger.NewZapLogger()
	defer logger.Sync()

	tasksRepo := mysql.NewTasksRepoMySQL(ctx, &cfg.MySQLDb)
	logger.Info("Tasks repo started successfully")

	usersRepo, client := mongo.NewUsersRepoMongoDB(ctx, &cfg.MongoDb)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	logger.Info("Users repo started successfully")

	sessionsRepo := redis.NewSessionsRepoRedis(ctx, &cfg.RedisDb)
	logger.Info("Sessions repo started successfully")

	usersService := service.NewUsersService(usersRepo)
	tasksService := service.NewTasksService(tasksRepo)
	sessionsService := service.NewSessionsService(sessionsRepo)

	mainService := service.NewService(*usersService, *sessionsService, *tasksService)

	httpHandler := httpHandler.NewHttpHandler(mainService, logger, tmpl)
	var server Server
	server = httpServer.NewHttpServer(ctx, httpHandler, &cfg.HTTPServer)

	if err := server.Run(ctx, logger); err != nil {
		logger.Fatal(ctx, err)
	}
}
