package handlers

import (
	"HW4/pkg/sessions"
	"HW4/pkg/users"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	ErrUserExist  = errors.New("user with this login exists")
	ErrBadId      = errors.New("bad id")
	ErrNoUser     = errors.New("no such user")
	ErrNoAuthData = errors.New("no username or password")
)

type SessionsHandler struct {
	SessionsRepo sessions.SessionsRepo
	UsersRepo    users.UsersRepo
	Logger       *zap.SugaredLogger
}

func (h *SessionsHandler) Login(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, ErrBadId.Error(), http.StatusBadGateway)
		h.Logger.Errorf(ErrBadId.Error())
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, ErrNoAuthData.Error(), http.StatusBadGateway)
		h.Logger.Errorf(ErrNoAuthData.Error())
		return
	}

	user := users.User{
		Id:       uint64(id),
		Login:    username,
		Password: password,
	}
	ok, err = h.UsersRepo.Check(r.Context(), &user)
	if err != nil {
		http.Error(w, ErrNoUser.Error(), http.StatusInternalServerError)
		h.Logger.Errorf(ErrNoUser.Error())
		return
	}
	if ok {
		http.Redirect(w, r, "/", http.StatusFound) // что-то в writer надо записать ????
	} else {
		http.Redirect(w, r, "/registration", http.StatusUnauthorized)
	}
}

func (h *SessionsHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusUnauthorized)
}

func (h *SessionsHandler) Registration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := users.User{
		Login:    vars["login"],
		Password: vars["password"],
	}

	err := h.UsersRepo.Add(r.Context(), &user)
	if err == ErrUserExist { // кастомно обработать
		renderJSON(w, err.Error())
		http.Redirect(w, r, "/registration", http.StatusUnauthorized)
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		renderJSON(w, "Registration success")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// renderJSON преобразует 'v' в формат JSON и записывает результат, в виде ответа, в w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	json, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
