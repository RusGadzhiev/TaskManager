package handlers

import (
	"HW4/pkg/tasks"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// надо ли писать код ошибки
// повторяющийся код, как и в репо таск

type TasksHandler struct {
	TasksRepo tasks.TasksRepo
	Logger    *zap.SugaredLogger
}

func (h *TasksHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.TasksRepo.GetAllTasks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	renderJSON(w, tasks)
}

func (h *TasksHandler) MyList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	tasks, err := h.TasksRepo.GetMyTasks(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	renderJSON(w, tasks)
}

func (h *TasksHandler) New(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	executor := vars["executor"]
	assign := false
	if executor != "" {
		assign = true
	}

	task := &tasks.Task{
		ID:          uint64(taskId),
		Owner:       vars["username"],
		Executor:    executor,
		Description: vars["description"],
		Completed:   false,
		Assigned:    assign,
	}
	err = h.TasksRepo.Add(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *TasksHandler) CreatedList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	tasks, err := h.TasksRepo.GetCreatedTasks(r.Context(), username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	renderJSON(w, tasks)
}

func (h *TasksHandler) Assign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["taskId"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	err = h.TasksRepo.Assign(r.Context(), uint64(taskId), vars["username"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TasksHandler) Unassign(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["taskId"])
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

	w.WriteHeader(http.StatusOK)
}

func (h *TasksHandler) Complete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskId, err := strconv.Atoi(vars["taskId"])
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

	w.WriteHeader(http.StatusOK)
}
