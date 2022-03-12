package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time := time.Now().Format("2006/01/02 15:04:05")
		method := r.Method
		host := r.URL.Host
		path := r.URL.Path
		log.SetFlags(0)
		log.Printf("%s %s %s%s\n", time, method, host, path)
		log.SetFlags(log.Llongfile | log.LstdFlags)

		handler.ServeHTTP(rw, r)
	})
}
