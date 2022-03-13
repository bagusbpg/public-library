package router

import (
	"net/http"
	_mw "plain-go/public-library/app/middleware"
	_user "plain-go/public-library/delivery/controller/user"
)

func Router(
	mux *http.ServeMux,
	user *_user.UserController,
) {
	mux.Handle("/login", (_mw.New(_mw.JSON, _mw.Logger).Then(user.Login())))
	mux.Handle("/users", (_mw.New(_mw.JSON, _mw.Logger).Then(user.SignUp())))
	// mux.Handle("/users/", (_mw.New(_mw.JSON, _mw.Logger, _mw.Authentication, _mw.ValidateId, _mw.Authorization).Then(user.GetUpdateDelete())))
}
