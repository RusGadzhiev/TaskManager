package handlers

import (
	"HW4/internal/tasks"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// можно ли обойтись без mux.vars , в мидлвея еще надо убрать
// сделать логику более реалистичной

type TasksHandler struct {
	TasksRepo tasks.TasksRepo
	Logger    *zap.SugaredLogger
	Tmpl      *template.Template
}

const (
	Description = "description"
	Executor    = "executor"
	UserName = "username"
)

func (h *TasksHandler) New(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, "create.html")
		return
	}

	vars := mux.Vars(r)
	executor := r.FormValue(Executor)
	assign := false
	if executor != "" {
		assign = true
	}

	task := &tasks.Task{
		Owner:       vars[UserName],
		Executor:    executor,
		Description: r.FormValue(Description),
		Completed:   false,
		Assigned:    assign,
	}
	taskId, err := h.TasksRepo.Add(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	renderJSON(w, taskId, h.Logger)
}

func (h *TasksHandler) List(w http.ResponseWriter, r *http.Request) {
	h.listSth(w, r, tasks.FilterAllTasks)
}

func (h *TasksHandler) MyList(w http.ResponseWriter, r *http.Request) {
	h.listSth(w, r, tasks.FilterMyTasks)
}

func (h *TasksHandler) CreatedList(w http.ResponseWriter, r *http.Request) {
	h.listSth(w, r, tasks.FilterCreatedTasks)
}

// сделать исполнителем можно только себя, можно забрать чуюжую задачу
func (h *TasksHandler) Assign(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, "assign.html")
		return
	}

	h.updateSth(w, r, tasks.FilterAssign)
}

// можно снять любую задачу, даже чужую
func (h *TasksHandler) Unassign(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, "unassign.html")
		return
	}

	h.updateSth(w, r, tasks.FilterUnassign)
}

// выполнить можно любую задачу, даже чужую
func (h *TasksHandler) Complete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.execTmpl(w, "complete.html")
		return
	}

	h.updateSth(w, r, tasks.FilterComplete)
}

func (h *TasksHandler) listSth(w http.ResponseWriter, r *http.Request, filter string) {

	vars := mux.Vars(r)
	username := vars[UserName]
	var tasksList []*tasks.Task
	var err error

	switch filter {
	case tasks.FilterAllTasks:
		tasksList, err = h.TasksRepo.GetAllTasks(r.Context())
	case tasks.FilterMyTasks:
		tasksList, err = h.TasksRepo.GetMyTasks(r.Context(), username)
	case tasks.FilterCreatedTasks:
		tasksList, err = h.TasksRepo.GetCreatedTasks(r.Context(), username)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}
	renderJSON(w, tasksList, h.Logger)
}

func (h *TasksHandler) updateSth(w http.ResponseWriter, r *http.Request, filter string) {
	taskId, err := strconv.Atoi(r.FormValue(tasks.TaskId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	vars := mux.Vars(r)
	switch filter {
		case tasks.FilterAssign:
			err = h.TasksRepo.Assign(r.Context(), uint64(taskId), vars[UserName])
		case tasks.FilterUnassign:
			err = h.TasksRepo.Unassign(r.Context(), uint64(taskId))
		case tasks.FilterComplete:
			err = h.TasksRepo.Complete(r.Context(), uint64(taskId))	
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}
}

// выпонлняет шаблон с именем tmpl, ответ в w записывает
func (h *TasksHandler) execTmpl(w http.ResponseWriter, tmpl string) {
	err := h.Tmpl.ExecuteTemplate(w, tmpl, nil)
	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}