package handlers

import (
	"HW4/pkg/tasks"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// можно ли обойтись без mux.vars , в мидлвея еще надо убрать
// сделать логику более реалистичной
// повторяющийся код, как и в репо таск

type TasksHandler struct {
	TasksRepo tasks.TasksRepo
	Logger    *zap.SugaredLogger
	Tmpl      *template.Template
}

const (
	Description = "description"
	TaskId      = "taskId"
	Executor    = "executor"
)

func (h *TasksHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TasksRepo.GetAllTasks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}
	h.Logger.Infof("URL: %s Status: %d", r.URL, http.StatusOK)
	renderJSON(w, tasks, h.Logger)
}

func (h *TasksHandler) MyList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars[UserName]
	tasks, err := h.TasksRepo.GetMyTasks(r.Context(), username)
	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.Logger.Infof("URL: %s Status: %d", r.URL, http.StatusOK)
	renderJSON(w, tasks, h.Logger)
}

func (h *TasksHandler) New(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
		if err != nil {
			h.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	h.Logger.Infof("URL: %s Status: %d", r.URL, http.StatusCreated)
	w.WriteHeader(http.StatusCreated)
	renderJSON(w, taskId, h.Logger)
}

func (h *TasksHandler) CreatedList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars[UserName]
	tasks, err := h.TasksRepo.GetCreatedTasks(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}
	h.Logger.Infof("URL: %s Status: %d", r.URL, http.StatusOK)
	renderJSON(w, tasks, h.Logger)
}

// сделать исполнителем можно только себя, можно забрать чуюжую задачу
func (h *TasksHandler) Assign(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		err := h.Tmpl.ExecuteTemplate(w, "assign.html", nil)
		if err != nil {
			h.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(r.FormValue(TaskId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	err = h.TasksRepo.Assign(r.Context(), uint64(taskId), vars[UserName])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}
	h.Logger.Infof("URL: %s Status: %d", r.URL, http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

// можно снять любую задачу, даже чужую
func (h *TasksHandler) Unassign(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		err := h.Tmpl.ExecuteTemplate(w, "unassign.html", nil)
		if err != nil {
			h.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	taskId, err := strconv.Atoi(r.FormValue(TaskId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	err = h.TasksRepo.Unassign(r.Context(), uint64(taskId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}
	h.Logger.Infof("URL: %s Status: %d", r.URL, http.StatusOK)
	w.WriteHeader(http.StatusOK)
}

// выполнить можно любую задачу, даже чужую
func (h *TasksHandler) Complete(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		err := h.Tmpl.ExecuteTemplate(w, "complete.html", nil)
		if err != nil {
			h.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	taskId, err := strconv.Atoi(r.FormValue(TaskId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	err = h.TasksRepo.Complete(r.Context(), uint64(taskId))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}
	h.Logger.Infof("URL: %s Status: %d", r.URL, http.StatusOK)
	w.WriteHeader(http.StatusOK)
}
