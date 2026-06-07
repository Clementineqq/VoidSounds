package middleware

import "net/http"

// использовать POST с _method=PUT/DELETE, шобы из формы редачить
func MethodOverride(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			method := r.URL.Query().Get("_method")
			if method == "PUT" || method == "DELETE" || method == "PATCH" {
				r.Method = method
			}
		}
		next.ServeHTTP(w, r)
	})
}
