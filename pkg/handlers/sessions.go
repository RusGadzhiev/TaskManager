package handlers

import (
	"HW4/pkg/sessions"
	"HW4/pkg/users"
	"encoding/json"
	"errors"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"
)

var (
	ErrBadId             = errors.New("bad id")
	ErrNoUser            = errors.New("no such user")
	ErrUserExist         = errors.New("user with this login exists")
	ErrIncorrectPassword = errors.New("incorrect password")
	letterRunes          = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

const (
	CookieName = "session_id"
)

type SessionsHandler struct {
	SessionsRepo sessions.SessionsRepo
	UsersRepo    users.UsersRepo
	Logger       *zap.SugaredLogger
	Tmpl         *template.Template
}

// предполагается что куки у пользователя нет
func (h *SessionsHandler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		err := h.Tmpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			h.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	user := users.User{
		UserName: r.FormValue(users.UserName),
		Password: r.FormValue(users.Password),
	}

	if !h.UsersRepo.IsUserExist(r.Context(), user.UserName) {
		h.Logger.Info(ErrNoUser.Error(), ": ", user.UserName)
		http.Error(w, ErrNoUser.Error(), http.StatusUnauthorized)
		return
	}

	pass, err := h.UsersRepo.GetPassword(r.Context(), user.UserName)
	if err != nil {
		h.Logger.Error(err.Error())
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}
	if pass != user.Password {
		h.Logger.Warn(ErrIncorrectPassword.Error())
		http.Error(w, ErrIncorrectPassword.Error(), http.StatusUnauthorized)
		return
	}

	cookieValue := randStringChars(32) // добавить уникальность куки
	cookie := http.Cookie{
		Name:    CookieName,
		Value:   cookieValue,
		Expires: time.Now().Add(72 * time.Hour),
	}
	dur := time.Until(cookie.Expires)
	err = h.SessionsRepo.Add(r.Context(), cookieValue, user.UserName, dur)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.Logger.Error(err.Error())
		return
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *SessionsHandler) Logout(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		err := h.Tmpl.ExecuteTemplate(w, "logout.html", nil)
		if err != nil {
			h.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	cookie, err := r.Cookie(CookieName)
	if err != nil {
		h.Logger.Info("New request",
			"method", r.Method,
			"remote_addr", r.RemoteAddr,
			"url", r.URL.Path,
			"error", "Permission denied")
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	h.SessionsRepo.Delete(r.Context(), cookie.Value)

	cookie.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/login", http.StatusUnauthorized)
}

// http://localhost:8080/registration
func (h *SessionsHandler) Registration(w http.ResponseWriter, r *http.Request) {

	h.Logger.Info("New request:",
		" method: ", r.Method,
		" remote_addr: ", r.RemoteAddr,
		" url: ", r.URL.Path)

	if r.Method == http.MethodGet {
		err := h.Tmpl.ExecuteTemplate(w, "registration.html", nil)
		if err != nil {
			h.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	user := users.User{
		UserName: r.FormValue(users.UserName),
		Password: r.FormValue(users.Password),
	}

	if h.UsersRepo.IsUserExist(r.Context(), user.UserName) {
		h.Logger.Info(ErrUserExist.Error())
		renderJSON(w, ErrUserExist.Error(), h.Logger)
		return
	}

	err := h.UsersRepo.Add(r.Context(), &user)
	if err != nil {
		h.Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		h.Logger.Info("Registration success")
		renderJSON(w, "Registration success", h.Logger)
	}
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

func randStringChars(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
