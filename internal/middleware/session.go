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
