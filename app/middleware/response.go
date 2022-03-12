package middleware

import "net/http"

func JSON(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		handler.ServeHTTP(rw, r)
	})
}
