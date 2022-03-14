package middleware

import (
	"log"
	"net/http"
	"strings"

	_model "plain-go/public-library/model"
)

func Authentication(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

		if token == "" {
			log.Println("missing or malformed jwt")
			_model.CreateResponse(rw, http.StatusBadRequest, "missing or malformed jwt", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}
