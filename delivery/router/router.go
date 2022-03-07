package router

import (
	"net/http"
	_user "plain-go/public-library/delivery/controller/user"
)

func Router(
	mux *http.ServeMux,
	userController *_user.UserController,
) {

	mux.Handle("/users", userController.CreateNewUser())
}
