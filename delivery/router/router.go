package router

import (
	"net/http"
	_user "plain-go/public-library/delivery/controller/user"
)

func Router(
	mux *http.ServeMux,
	user *_user.UserController,
) {
	mux.Handle("/login", user.Login())
	mux.Handle("/users", user.SignUp())
	mux.Handle("/users/", user.GetUpdateDelete())
}
