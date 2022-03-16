package middleware

import (
	"log"
	"net/http"
	_model "plain-go/public-library/model"
	"strconv"
	"strings"
)

func ValidateId(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		str := strings.SplitAfter(r.URL.Path, "/")
		_, err := strconv.Atoi(str[len(str)-1])

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "invalid id", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}
