package middleware

import (
	"log"
	"net/http"
	_model "plain-go/public-library/model"
	"strconv"
)

func ValidateId(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, err := strconv.Atoi(GetParam(r, 0))

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "invalid id", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}
