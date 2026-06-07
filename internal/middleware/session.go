package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitSessionStore(secret string) {
	Store = sessions.NewCookieStore([]byte(secret))
}

type contextKey string

const UserContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Store == nil {
			next.ServeHTTP(w, r)
			return
		}

		session, err := Store.Get(r, "user-session")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		userID, ok := session.Values["user_id"].(int)
		if ok && userID != 0 {
			ctx := context.WithValue(r.Context(), UserContextKey, userID)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) int {
	if userID, ok := r.Context().Value(UserContextKey).(int); ok {
		return userID
	}
	return 0
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if GetUserID(r) == 0 {
			// HTMX-запрос или обычный переход в браузере
			if r.Header.Get("HX-Request") == "true" {
				// для HTMX: отправляем спец-заголовок для редиректа
				w.Header().Set("HX-Redirect", "/login")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// для браузера стандартный HTTP-редирект
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := Store.Get(r, "user-session")
			userRole, _ := session.Values["user_role"].(string)

			if userRole != role {
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/")
					w.WriteHeader(http.StatusForbidden)
					w.Write([]byte("Доступ запрещён"))
					return
				}
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
