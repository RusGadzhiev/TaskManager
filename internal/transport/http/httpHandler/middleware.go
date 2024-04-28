package httpHandler

import (
	"net/http"

	"github.com/RusGadzhiev/TaskManager/internal/service"
	"github.com/gorilla/mux"
)

func (h *HttpHandler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.logger.Info("New request: ", "method - ", r.Method, " remote_addr - ", r.RemoteAddr, " url - ", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// по значению куки устанавливает значение username
func (h *HttpHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(service.CookieName)
		if err != nil {
			h.logger.Info("Permission denied")
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}

		username, err := h.service.GetUserByCookie(r.Context(), cookie.Value)
		if err == service.ErrNoUserBySession {
			h.logger.Info(err.Error())
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		} else if err != nil {
			h.logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		h.logger.Info("Auth Success")
		mux.Vars(r)[service.UserName] = username
		next.ServeHTTP(w, r)
	})
}

func (h *HttpHandler) PanicRecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				h.logger.Error("Url: ", r.URL.Path, " Recovered: ", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
