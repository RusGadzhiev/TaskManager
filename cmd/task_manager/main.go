package main

import (
	"HW4/internal/config"
	"HW4/pkg/handlers"
	"HW4/pkg/sessions"
	"HW4/pkg/tasks"
	"HW4/pkg/users"
	"context"
	"net/http"

	"HW4/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {

	cfg := config.MustLoad()

	ctx := context.Background() // сделай нормальный контекст

	logger := logger.NewZapLogger()
	defer logger.Sync()

	tasksRepo := tasks.NewTasksRepoMySQL(ctx, &cfg.MySQLDb)

	usersRepo, client := users.NewUsersRepoMongoDB(ctx, &cfg.MongoDb)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	sessionsRepo := sessions.NewSessionsRepoRedis(ctx, &cfg.RedisDb)

	// var _ *tasks.TasksRepo = (*tasks.TasksRepoMySQL)(nil) // почему не работает проверка

	tasksHandler := &handlers.TasksHandler{
		TasksRepo: tasksRepo,
		Logger:    logger,
	}

	sessionsHandler := &handlers.SessionsHandler{
		SessionsRepo: sessionsRepo,
		UsersRepo:    usersRepo,
		Logger:       logger,
	}

	/*middleware := &handlers.MiddlewareHandler{
		SessionsRepo: sessionsRepo,
		UsersRepo:    usersRepo,
		Logger:       logger,
	}*/

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", tasksHandler.List).Methods("GET")
	r.HandleFunc("/login", sessionsHandler.Login).Methods("POST")
	r.HandleFunc("/logout", sessionsHandler.Logout).Methods("POST")
	r.HandleFunc("/registration", sessionsHandler.Registration).Methods("POST")

	taskRouter := mux.NewRouter()
	taskRouter.StrictSlash(true)
	taskRouter.HandleFunc("/tasks", tasksHandler.MyList).Methods("GET")
	taskRouter.HandleFunc("/tasks/created", tasksHandler.CreatedList).Methods("GET")
	taskRouter.HandleFunc("/tasks/new", tasksHandler.New).Methods("POST")
	taskRouter.HandleFunc("/tasks/assign/{id}", tasksHandler.Assign).Methods("POST")
	taskRouter.HandleFunc("/tasks/unassign/{id}", tasksHandler.Unassign).Methods("POST")
	taskRouter.HandleFunc("/tasks/complete/{id}", tasksHandler.Complete).Methods("POST")

	/*taskRouter.Use(func(h http.Handler) http.Handler {
		return middleware.AuthMiddleware(h)
	})*/
	// подключи другие middleware
	http.ListenAndServe(":8080", r) // тут лучше ListenAndServeTLS
}
