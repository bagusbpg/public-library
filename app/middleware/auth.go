package middleware

import (
	"log"
	"net/http"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
	"strconv"
	"strings"
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

func AuthorizedById(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

		loginId, _, _ := _helper.ExtractToken(token)

		userId, _ := strconv.Atoi(GetParam(r)[0])

		if loginId != userId {
			log.Println("forbidden")
			_model.CreateResponse(rw, http.StatusForbidden, "forbidden", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}

func MemberOnlyAuthorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

		_, role, err := _helper.ExtractToken(token)

		if err != nil {
			_model.CreateResponse(rw, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		if role != "member" {
			log.Println("forbidden")
			_model.CreateResponse(rw, http.StatusForbidden, "forbidden", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}

func LibrarianOnlyAuthorization(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

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
