package middleware

import (
	"log"
	"net/http"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
	"strconv"
	"strings"
)

func MemberOnlydAuthorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

		if token == "" {
			log.Println("missing or malformed jwt")
			_model.CreateResponse(rw, http.StatusBadRequest, "missing or malformed jwt", nil)
			return
		}

		loginId, _, err := _helper.ExtractToken(token)

		if err != nil {
			_model.CreateResponse(rw, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		str := strings.SplitAfter(r.URL.Path, "/")
		userId, _ := strconv.Atoi(str[len(str)-1])

		if loginId != userId {
			log.Println("forbidden")
			_model.CreateResponse(rw, http.StatusForbidden, "forbidden", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}

func MemberAndLibrarianAuthorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

		if token == "" {
			log.Println("missing or malformed jwt")
			_model.CreateResponse(rw, http.StatusBadRequest, "missing or malformed jwt", nil)
			return
		}

		loginId, role, err := _helper.ExtractToken(token)

		if err != nil {
			_model.CreateResponse(rw, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		str := strings.SplitAfter(r.URL.Path, "/")
		userId, _ := strconv.Atoi(str[len(str)-1])

		if loginId != userId && role != "Librarian" {
			log.Println("forbidden")
			_model.CreateResponse(rw, http.StatusForbidden, "forbidden", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}

func LibrarianAuthorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

		if token == "" {
			log.Println("missing or malformed jwt")
			_model.CreateResponse(rw, http.StatusBadRequest, "missing or malformed jwt", nil)
			return
		}

		_, role, err := _helper.ExtractToken(token)

		if err != nil {
			_model.CreateResponse(rw, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		if role != "Librarian" {
			log.Println("forbidden")
			_model.CreateResponse(rw, http.StatusForbidden, "forbidden", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}
