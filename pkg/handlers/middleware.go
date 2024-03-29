package handlers

import (
	"HW4/pkg/sessions"
	"HW4/pkg/users"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// наверно структура не нужна
type MiddlewareHandler struct {
	SessionsRepo sessions.SessionsRepo
	UsersRepo    users.UsersRepo
	Logger       *zap.SugaredLogger
}

func (h *MiddlewareHandler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		h.Logger.Info("New request", "method", r.Method, "remote_addr", r.RemoteAddr, "url", r.URL.Path, "time", time.Since(start))
	})
}

func (h *MiddlewareHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session")
		if err != nil {
			h.Logger.Info("New request",
				"method", r.Method,
				"remote_addr", r.RemoteAddr,
				"url", r.URL.Path,
				"error", "Permission denied")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, ErrBadId.Error(), http.StatusBadGateway)
			return
		}
		s, err := h.SessionsRepo.GetSession(r.Context(), uint64(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if s.Cookie.Value != session.Value {
			h.Logger.Info("New request",
				"method", r.Method,
				"remote_addr", r.RemoteAddr,
				"url", r.URL.Path,
				"error", "Cookies dont same")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *MiddlewareHandler) PanicRecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				h.Logger.Info("New request",
					"method", r.Method,
					"remote_addr", r.RemoteAddr,
					"url", r.URL.Path,
					"recovered", err)
				http.Error(w, "Internal server error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
