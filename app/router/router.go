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
	_wish "plain-go/public-library/controller/wish"
	_model "plain-go/public-library/model"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

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
	wish *_wish.WishController,
) http.HandlerFunc {
	routes := []route{
		NewRoute(http.MethodPost, "/login", _mw.Do(_mw.JSONRequest).Then(user.Login()).ServeHTTP),
		NewRoute(http.MethodPost, "/users", _mw.Do(_mw.JSONRequest).Then(user.SignUp()).ServeHTTP),
		NewRoute(http.MethodGet, "/users", _mw.Do(_mw.Authentication, _mw.LibrarianOnlyAuthorization).Then(user.GetAll()).ServeHTTP),
		NewRoute(http.MethodGet, "/users/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.AuthorizedById).Then(user.Get()).ServeHTTP),
		NewRoute(http.MethodPut, "/users/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.AuthorizedById, _mw.JSONRequest).Then(user.Update()).ServeHTTP),
		NewRoute(http.MethodDelete, "/users/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.AuthorizedById).Then(user.Delete()).ServeHTTP),
		NewRoute(http.MethodPost, "/books", _mw.Do(_mw.JSONRequest, _mw.LibrarianOnlyAuthorization).Then(book.Create()).ServeHTTP),
		NewRoute(http.MethodGet, "/books", book.GetAll().ServeHTTP),
		NewRoute(http.MethodGet, "/books/(.+)", _mw.Do(_mw.ValidateId).Then(book.Get()).ServeHTTP),
		NewRoute(http.MethodPut, "/books/(.+)", _mw.Do(_mw.ValidateId, _mw.JSONRequest, _mw.Authentication, _mw.LibrarianOnlyAuthorization).Then(book.Update()).ServeHTTP),
		NewRoute(http.MethodDelete, "/books/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.LibrarianOnlyAuthorization).Then(book.Delete()).ServeHTTP),
		NewRoute(http.MethodPost, "/favorites/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.AuthorizedById, _mw.JSONRequest).Then(favorite.AddBook()).ServeHTTP),
		NewRoute(http.MethodDelete, "/favorites/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.AuthorizedById, _mw.JSONRequest).Then(favorite.RemoveBook()).ServeHTTP),
		NewRoute(http.MethodGet, "/favorites/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.AuthorizedById).Then(favorite.GetAll()).ServeHTTP),
		NewRoute(http.MethodGet, "/wishes", _mw.Do(_mw.Authentication, _mw.LibrarianOnlyAuthorization).Then(wish.GetAll()).ServeHTTP),
		NewRoute(http.MethodPost, "/wishes/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.AuthorizedById, _mw.JSONRequest).Then(wish.AddBook()).ServeHTTP),
		NewRoute(http.MethodGet, "/wishes/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.AuthorizedById, _mw.MemberOnlyAuthorization).Then(wish.GetAllByUser()).ServeHTTP),
		NewRoute(http.MethodPut, "/wishes/(.+)/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.AuthorizedById, _mw.MemberOnlyAuthorization).Then(wish.Update()).ServeHTTP),
		NewRoute(http.MethodDelete, "/wishes/(.+)/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication, _mw.AuthorizedById, _mw.MemberOnlyAuthorization).Then(wish.RemoveBook()).ServeHTTP),
		// NewRoute(http.MethodGet, "/reviews", _mw.Do(_mw.Authentication, _mw.LibrarianOnlyAuthorization).Then(review.GetAll()).ServeHTTP),
		// NewRoute(http.MethodPost, "/reviews/(.+)", _mw.Do(_mw.JSONRequest, _mw.ValidateId, _mw.Authentication).Then(review.Create()).ServeHTTP),
		// NewRoute(http.MethodGet, "/reviews/(.+)", _mw.Do(_mw.ValidateId, _mw.Authentication).Then(review.Get()).ServeHTTP),
		// NewRoute(http.MethodPut, "/reviews/(.+)", _mw.Do(_mw.ValidateId, _mw.JSONRequest, _mw.Authentication))
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

		allowed := []string{}

		for _, route := range routes {
			matches := route.regex.FindStringSubmatch(r.URL.Path)

			// if route is matched
			if len(matches) > 0 {
				// if method is not matched, continue
				if r.Method != route.method {
					allowed = append(allowed, route.method)
					continue
				}

				// if method is matched, call handler
				// embed params to context
				ctx := context.WithValue(r.Context(), _mw.ContextKey{}, matches[1:])
				route.handler(rw, r.WithContext(ctx))
				return
			}
		}

		// if endpoint is matched, but method is not
		if len(allowed) > 0 {
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		// if endpoint is not matched
		log.Println("endpoint not found")
		_model.CreateResponse(rw, http.StatusNotFound, "endpoint not found", nil)
	}
}
