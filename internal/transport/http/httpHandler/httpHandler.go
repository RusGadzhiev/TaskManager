package httpHandler

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"github.com/RusGadzhiev/TaskManager/internal/service"
)

// куда закинуть функции утилиты
// добавь контекст к каждому запросу
// попробуй вынести контекст в middleware
// время для контекста должно браться из конфига

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
	// проверяет пароль
	CheckPass(ctx context.Context, username string)
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

// ни от чего не зависит, использует только свои модели, легко переиспользуется
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

	if r.Method == http.MethodGet {
		h.execTmpl(w, "create.html")
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
	taskId, err := h.service.Add(r.Context(), task)
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

// сделать исполнителем можно только себя, можно забрать чуюжую задачу
func (h *HttpHandler) Assign(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, "assign.html")
		return
	}

	h.updateSth(w, r, service.FilterAssign)
}

// можно снять любую задачу, даже чужую
func (h *HttpHandler) Unassign(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, "unassign.html")
		return
	}

	h.updateSth(w, r, service.FilterUnassign)
}

// выполнить можно любую задачу, даже чужую
func (h *HttpHandler) Complete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, "complete.html")
		return
	}

	h.updateSth(w, r, service.FilterComplete)
}

// предполагается что куки у пользователя нет
func (h *HttpHandler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, "login.html")
		return
	}

	newUser := service.User{
		UserName: r.FormValue(service.UserName),
		Password: r.FormValue(service.Password),
	}

	err := h.service.Authentificate(r.Context(), &newUser)
	if err == service.ErrNoUser ||  err == service.ErrIncorrectPassword {
		h.logger.Info(err.Error(), " user: ", newUser.UserName)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else if err != nil {
		h.logger.Error(err.Error())
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}

	session, err := h.service.AddCookie(r.Context(), newUser.UserName)
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

	if r.Method == http.MethodGet {
		h.execTmpl(w, "logout.html")
		return
	}

	cookie, err := r.Cookie(service.CookieName)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.DeleteCookie(r.Context(), cookie.Value)
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

	if r.Method == http.MethodGet {
		h.execTmpl(w, "registration.html")
		return
	}

	user := service.User{
		UserName: r.FormValue(service.UserName),
		Password: r.FormValue(service.Password),
	}

	err := h.service.AddUser(r.Context(), &user)
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

	vars := mux.Vars(r)
	username := vars[service.UserName]
	var tasksList []*service.Task
	var err error

	switch filter {
	case service.FilterAllTasks:
		tasksList, err = h.service.GetAllTasks(r.Context())
	case service.FilterMyTasks:
		tasksList, err = h.service.GetMyTasks(r.Context(), username)
	case service.FilterCreatedTasks:
		tasksList, err = h.service.GetCreatedTasks(r.Context(), username)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}

	renderJSON(w, tasksList, h.logger)
}

func (h *HttpHandler) updateSth(w http.ResponseWriter, r *http.Request, filter string) {
	taskId, err := strconv.Atoi(r.FormValue(service.TaskId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.logger.Error(err.Error())
		return
	}
	vars := mux.Vars(r)
	switch filter {
	case service.FilterAssign:
		err = h.service.Assign(r.Context(), uint64(taskId), vars[service.UserName])
	case service.FilterUnassign:
		err = h.service.Unassign(r.Context(), uint64(taskId))
	case service.FilterComplete:
		err = h.service.Complete(r.Context(), uint64(taskId))
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
