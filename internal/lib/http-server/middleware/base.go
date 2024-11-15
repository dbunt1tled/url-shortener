package middleware

import (
	"net/http"
)

func Base() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			headers := w.Header()
			headers.Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}

}
