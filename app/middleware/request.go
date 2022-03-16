package middleware

import (
	"log"
	"net/http"
	_model "plain-go/public-library/model"
)

func JSONRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if contentType := r.Header.Get("content-type"); contentType != "application/json" {
			log.Println("unsupported content type")
			_model.CreateResponse(rw, http.StatusUnsupportedMediaType, "unsupported content type", nil)
			return
		}

		handler.ServeHTTP(rw, r)
	})
}
