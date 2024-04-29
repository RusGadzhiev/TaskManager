package httpHandler

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/RusGadzhiev/TaskManager/internal/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	templateCreate       = "create.html"
	templateAssign       = "assign.html"
	templateUnassign     = "unassign.html"
	templateRegistration = "registration.html"
	templateLogout       = "logout.html"
	templateComplete     = "complete.html"
	templateLogin        = "login.html"
)

type TasksService interface {
	GetAllTasks(ctx context.Context) ([]*service.Task, error)
	GetCreatedTasks(ctx context.Context, username string) ([]*service.Task, error)
	GetMyTasks(ctx context.Context, username string) ([]*service.Task, error)
	// возвращает id вставленной задачи
	Add(ctx context.Context, task *service.Task) (uint64, error)
	Assign(ctx context.Context, taskId uint64, username string) error
	Unassign(ctx context.Context, taskId uint64) error
	Complete(ctx context.Context, taskId uint64) error
}

type UsersService interface {
	// возвращает ошибку service.ErrNoUser если юзера нет, ошбику service.ErrIncorrectPassword если пароль неверный
	Authentificate(ctx context.Context, user *service.User) error
	// возвращает ошибку service.ErrUserExist если юзер уже есть
	AddUser(ctx context.Context, user *service.User) error
}

type SessionsService interface {
	DeleteCookie(ctx context.Context, cookieVal string) error
	// возвращает значение созданной куки и продолжительность действия
	AddCookie(ctx context.Context, username string) (*service.Session, error)
	// возвращает юзернейм по значению куки
	GetUserByCookie(ctx context.Context, cookieVal string) (string, error)
}

type Service interface {
	UsersService
	TasksService
	SessionsService
}

type HttpHandler struct {
	service Service
	tmpl    *template.Template
	logger  *zap.SugaredLogger
}

func (h *HttpHandler) New(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10 * time.Second)
	defer cancel()

	if r.Method == http.MethodGet {
		h.execTmpl(w, templateCreate)
		return
	}

	vars := mux.Vars(r)
	executor := r.FormValue(service.Executor)
	assign := false
	if executor != "" {
		assign = true
	}

	task := &service.Task{
		Owner:       vars[service.UserName],
		Executor:    executor,
		Description: r.FormValue(service.Description),
		Completed:   false,
		Assigned:    assign,
	}
	taskId, err := h.service.Add(ctx, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	renderJSON(w, taskId, h.logger)
}

func (h *HttpHandler) List(w http.ResponseWriter, r *http.Request) {
	h.listSth(w, r, service.FilterAllTasks)
}

func (h *HttpHandler) MyList(w http.ResponseWriter, r *http.Request) {
	h.listSth(w, r, service.FilterMyTasks)
}

func (h *HttpHandler) CreatedList(w http.ResponseWriter, r *http.Request) {
	h.listSth(w, r, service.FilterCreatedTasks)
}

func (h *HttpHandler) Assign(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, templateAssign)
		return
	}

	h.updateSth(w, r, service.FilterAssign)
}

func (h *HttpHandler) Unassign(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, templateUnassign)
		return
	}

	h.updateSth(w, r, service.FilterUnassign)
}

func (h *HttpHandler) Complete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, templateComplete)
		return
	}

	h.updateSth(w, r, service.FilterComplete)
}

func (h *HttpHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20 * time.Second)
	defer cancel()

	if r.Method == http.MethodGet {
		h.execTmpl(w, templateLogin)
		return
	}

	newUser := service.User{
		UserName: r.FormValue(service.UserName),
		Password: r.FormValue(service.Password),
	}

	err := h.service.Authentificate(ctx, &newUser)
	if err == service.ErrNoUser || err == service.ErrIncorrectPassword {
		h.logger.Info(err.Error(), " user: ", newUser.UserName)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}

	session, err := h.service.AddCookie(ctx, newUser.UserName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}
	cookie := http.Cookie{
		Name:    service.CookieName,
		Value:   session.CookieVal,
		Expires: time.Now().Add(session.Dur),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *HttpHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10 * time.Second)
	defer cancel()

	if r.Method == http.MethodGet {
		h.execTmpl(w, templateLogout)
		return
	}

	cookie, err := r.Cookie(service.CookieName)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.DeleteCookie(ctx, cookie.Value)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusUnauthorized)
}

func (h *HttpHandler) Registration(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10 * time.Second)
	defer cancel()

	if r.Method == http.MethodGet {
		h.execTmpl(w, templateRegistration)
		return
	}

	user := service.User{
		UserName: r.FormValue(service.UserName),
		Password: r.FormValue(service.Password),
	}

	err := h.service.AddUser(ctx, &user)
	if err == service.ErrUserExist {
		h.logger.Info(service.ErrUserExist.Error(), ": ", user.UserName)
		http.Error(w, service.ErrUserExist.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Info("Registration success")
	renderJSON(w, "Registration success", h.logger)
}

// renderJSON преобразует 'v' в формат JSON и записывает результат, в виде ответа, в w.
func renderJSON(w http.ResponseWriter, v interface{}, logger *zap.SugaredLogger) {
	json, err := json.Marshal(v)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func (h *HttpHandler) listSth(w http.ResponseWriter, r *http.Request, filter string) {
	ctx, cancel := context.WithTimeout(r.Context(), 10 * time.Second)
	defer cancel()

	vars := mux.Vars(r)
	username := vars[service.UserName]
	var tasksList []*service.Task
	var err error

	switch filter {
	case service.FilterAllTasks:
		tasksList, err = h.service.GetAllTasks(ctx)
	case service.FilterMyTasks:
		tasksList, err = h.service.GetMyTasks(ctx, username)
	case service.FilterCreatedTasks:
		tasksList, err = h.service.GetCreatedTasks(ctx, username)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}

	renderJSON(w, tasksList, h.logger)
}

func (h *HttpHandler) updateSth(w http.ResponseWriter, r *http.Request, filter string) {
	ctx, cancel := context.WithTimeout(r.Context(), 10 * time.Second)
	defer cancel()

	taskId, err := strconv.Atoi(r.FormValue(service.TaskId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}
	vars := mux.Vars(r)
	switch filter {
	case service.FilterAssign:
		err = h.service.Assign(ctx, uint64(taskId), vars[service.UserName])
	case service.FilterUnassign:
		err = h.service.Unassign(ctx, uint64(taskId))
	case service.FilterComplete:
		err = h.service.Complete(ctx, uint64(taskId))
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}
}

// выпонлняет шаблон с именем tmpl, ответ в w записывает
func (h *HttpHandler) execTmpl(w http.ResponseWriter, tmpl string) {
	err := h.tmpl.ExecuteTemplate(w, tmpl, nil)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *HttpHandler) Router() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", h.List).Methods("GET")
	r.HandleFunc("/login", h.Login).Methods("POST", "GET")
	r.HandleFunc("/logout", h.Logout).Methods("POST", "GET")
	r.HandleFunc("/registration", h.Registration).Methods("POST", "GET")
	r.Handle("/tasks", h.AuthMiddleware(http.HandlerFunc(h.MyList))).Methods("GET")
	r.Handle("/tasks/created", h.AuthMiddleware(http.HandlerFunc(h.CreatedList))).Methods("GET")
	r.Handle("/tasks/new", h.AuthMiddleware(http.HandlerFunc(h.New))).Methods("POST", "GET")
	r.Handle("/tasks/assign", h.AuthMiddleware(http.HandlerFunc(h.Assign))).Methods("POST", "GET")
	r.Handle("/tasks/unassign", h.AuthMiddleware(http.HandlerFunc(h.Unassign))).Methods("POST", "GET")
	r.Handle("/tasks/complete", h.AuthMiddleware(http.HandlerFunc(h.Complete))).Methods("POST", "GET")

	r.Use(func(hdl http.Handler) http.Handler {
		return h.PanicRecoverMiddleware(hdl)
	})
	r.Use(func(hdl http.Handler) http.Handler {
		return h.LoggingMiddleware(hdl)
	})
	return r
}

func NewHttpHandler(s Service, logger *zap.SugaredLogger, tmpl *template.Template) *HttpHandler {
	return &HttpHandler{
		service: s,
		logger:  logger,
		tmpl:    tmpl,
	}
}
