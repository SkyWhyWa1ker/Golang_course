package middleware

import "net/http"

func APIKey(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("X-API-KEY") != "secret12345" {
			http.Error(w, `{"error":"unauthorized"}`, 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}
