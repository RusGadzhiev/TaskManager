package handlers

import (
	"github.com/RusGadzhiev/TaskManager/internal/sessions"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *SessionsHandler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Logger.Info("New request: ", "method - ", r.Method, " remote_addr - ", r.RemoteAddr, " url - ", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// по значению куки устанавливает значение username
func (h *SessionsHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(CookieName)
		if err != nil {
			h.Logger.Info("Permission denied")
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		username, err := h.SessionsRepo.GetUser(r.Context(), cookie.Value)
		if err == sessions.ErrNoUserBySession {
			h.Logger.Info(err.Error())
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		} else if err != nil {
			h.Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		h.Logger.Info("Auth Success")
		mux.Vars(r)[UserName] = username
		next.ServeHTTP(w, r)
	})
}

func (h *SessionsHandler) PanicRecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				h.Logger.Error("Url: ", r.URL.Path, " Recovered: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
