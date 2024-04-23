package main

import (
	"HW4/internal/config"
	"HW4/internal/handlers"
	"HW4/internal/sessions"
	"HW4/internal/tasks"
	"HW4/internal/users"
	"context"
	"html/template"
	"net/http"

	"HW4/pkg/logger"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.MustLoad()
	tmpl := template.Must(template.ParseGlob("./templates/*"))

	ctx := context.Background() // сделай нормальный контекст

	logger := logger.NewZapLogger()
	defer logger.Sync()

	tasksRepo := tasks.NewTasksRepoMySQL(ctx, &cfg.MySQLDb)
	logger.Info("Tasks repo started successfully")

	usersRepo, client := users.NewUsersRepoMongoDB(ctx, &cfg.MongoDb)
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	logger.Info("Users repo started successfully")

	sessionsRepo := sessions.NewSessionsRepoRedis(ctx, &cfg.RedisDb)
	logger.Info("Sessions repo started successfully")
	// var _ *tasks.TasksRepo = (*tasks.TasksRepoMySQL)(nil) // почему не работает проверка

	tasksHandler := &handlers.TasksHandler{
		TasksRepo: tasksRepo,
		Logger:    logger,
		Tmpl:      tmpl,
	}

	sessionsHandler := &handlers.SessionsHandler{
		SessionsRepo: sessionsRepo,
		UsersRepo:    usersRepo,
		Logger:       logger,
		Tmpl:         tmpl,
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	// повесь и другие методы
	r.HandleFunc("/", tasksHandler.List).Methods("GET")
	r.HandleFunc("/login", sessionsHandler.Login).Methods("POST", "GET")
	r.HandleFunc("/logout", sessionsHandler.Logout).Methods("POST", "GET")
	r.HandleFunc("/registration", sessionsHandler.Registration).Methods("POST", "GET")
	r.Handle("/tasks", sessionsHandler.AuthMiddleware(http.HandlerFunc(tasksHandler.MyList))).Methods("GET")
	r.Handle("/tasks/created", sessionsHandler.AuthMiddleware(http.HandlerFunc(tasksHandler.CreatedList))).Methods("GET")
	r.Handle("/tasks/new", sessionsHandler.AuthMiddleware(http.HandlerFunc(tasksHandler.New))).Methods("POST", "GET")
	r.Handle("/tasks/assign", sessionsHandler.AuthMiddleware(http.HandlerFunc(tasksHandler.Assign))).Methods("POST", "GET")
	r.Handle("/tasks/unassign", sessionsHandler.AuthMiddleware(http.HandlerFunc(tasksHandler.Unassign))).Methods("POST", "GET")
	r.Handle("/tasks/complete", sessionsHandler.AuthMiddleware(http.HandlerFunc(tasksHandler.Complete))).Methods("POST", "GET")

	r.Use(func(h http.Handler) http.Handler {
		return sessionsHandler.PanicRecoverMiddleware(h)
	})
	r.Use(func(h http.Handler) http.Handler {
		return sessionsHandler.LoggingMiddleware(h)
	})
	logger.Info("Start listen at " + cfg.HTTPServer.Host + ":" + cfg.HTTPServer.Port)
	http.ListenAndServe(":"+cfg.HTTPServer.Port, r) // тут лучше ListenAndServeTLS
}
