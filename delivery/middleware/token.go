package middleware

import (
	"log"
	"net/http"
	"strings"

	_common "plain-go/public-library/delivery/common"
)

func Authorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

		if token == "" {
			log.Println("missing or malformed jwt")
			_common.CreateResponse(rw, http.StatusBadRequest, "missing or malformed jwt", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}
