package router

import (
	"net/http"
	_user "plain-go/public-library/delivery/controller/user"
	_mw "plain-go/public-library/delivery/middleware"
)

func Router(
	mux *http.ServeMux,
	user *_user.UserController,
) {
	mux.Handle("/login", (_mw.New(_mw.JSON, _mw.Logger).Then(user.Login())))
	mux.Handle("/users", (_mw.New(_mw.JSON, _mw.Logger).Then(user.SignUp())))
	mux.Handle("/users/", (_mw.New(_mw.JSON, _mw.Logger, _mw.Authorization).Then(user.GetUpdateDelete())))
}
