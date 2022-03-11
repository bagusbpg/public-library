package router

import (
	"net/http"
	_user "plain-go/public-library/delivery/controller/user"
	_middleware "plain-go/public-library/delivery/middleware"
)

func Router(
	mux *http.ServeMux,
	user *_user.UserController,
) {
	mux.Handle("/login", _middleware.JSONResponse(_middleware.Logger(user.Login())))
	mux.Handle("/users", _middleware.JSONResponse(_middleware.Logger(user.SignUp())))
	mux.Handle("/users/", _middleware.JSONResponse(_middleware.Authorization(_middleware.Logger(user.GetUpdateDelete()))))
}
