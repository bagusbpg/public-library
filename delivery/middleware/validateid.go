package middleware

import (
	"log"
	"net/http"
	_common "plain-go/public-library/delivery/common"
	"strconv"
	"strings"
)

func ValidateId(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, err := strconv.Atoi(strings.SplitAfter(r.URL.Path, "/")[2])

		if err != nil {
			log.Println(err)
			_common.CreateResponse(rw, http.StatusBadRequest, "invalid user id", nil)
			return
		}
	})
}
