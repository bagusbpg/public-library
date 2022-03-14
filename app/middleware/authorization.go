package middleware

import (
	"log"
	"net/http"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
	"strconv"
	"strings"
)

func Authorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

		loginId, role, err := _helper.ExtractToken(token)

		if err != nil {
			_model.CreateResponse(rw, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		userId, _ := strconv.Atoi(strings.SplitAfter(r.URL.Path, "/")[2])

		if loginId != userId && role != "Administrator" {
			log.Println("forbidden")
			_model.CreateResponse(rw, http.StatusForbidden, "forbidden", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}
