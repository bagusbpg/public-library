package middleware

import (
	"log"
	"net/http"
	_model "plain-go/public-library/model"
)

func POST(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}

func GET(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}

func PUT(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}

func DELETE(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}
