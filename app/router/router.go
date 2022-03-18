package router

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"time"

	_mw "plain-go/public-library/app/middleware"
	_book "plain-go/public-library/controller/book"
	_favorite "plain-go/public-library/controller/favorite"
	_user "plain-go/public-library/controller/user"
	_model "plain-go/public-library/model"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

type ctxKey struct{}

func NewRoute(method string, pattern string, handler http.HandlerFunc) route {
	return route{
		method,
		regexp.MustCompile("^" + pattern + "$"),
		handler,
	}
}

func Router(
	user *_user.UserController,
	book *_book.BookController,
	favorite *_favorite.FavoriteController,
) http.HandlerFunc {
	routes := []route{
		NewRoute(http.MethodPost, "/login", _mw.New(_mw.JSONRequest).Then(user.Login()).ServeHTTP),
		NewRoute(http.MethodPost, "/users", _mw.New(_mw.JSONRequest).Then(user.SignUp()).ServeHTTP),
		NewRoute(http.MethodGet, "/users", _mw.New(_mw.LibrarianAuthorization).Then(user.GetAll()).ServeHTTP),
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		// set JSON format as response
		rw.Header().Add("Content-Type", "application/json")

		// set logger
		time := time.Now().Format("2006/01/02 15:04:05")
		method := r.Method
		host := r.URL.Host
		path := r.URL.Path
		log.SetFlags(0)
		log.Printf("%s %s %s%s\n", time, method, host, path)
		log.SetFlags(log.Llongfile | log.LstdFlags)

		notAllowed := []string{}

		for _, route := range routes {
			matches := route.regex.FindStringSubmatch(r.URL.Path)

			if len(matches) > 0 {
				if r.Method != route.method {
					notAllowed = append(notAllowed, route.method)
					continue
				}

				ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])

				route.handler(rw, r.WithContext(ctx))
				return
			}
		}

		if len(notAllowed) > 0 {
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		log.Println("endpoint not found")
		_model.CreateResponse(rw, http.StatusNotFound, "endpoint not found", nil)
	}
}
