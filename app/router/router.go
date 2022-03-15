package router

import (
	"net/http"
	_mw "plain-go/public-library/app/middleware"
	_book "plain-go/public-library/controller/book"
	_user "plain-go/public-library/controller/user"
)

func Router(
	mux *http.ServeMux,
	user *_user.UserController,
	book *_book.BookController,
) {
	mux.Handle("/login", (_mw.New(_mw.JSONResponse, _mw.Logger).Then(user.Login())))
	mux.Handle("/users", (_mw.New(_mw.JSONResponse, _mw.Logger).Then(user.SignUp())))
	mux.Handle("/users/", (_mw.New(_mw.JSONResponse, _mw.Logger, _mw.Authentication, _mw.ValidateId, _mw.Authorization).Then(user.GetUpdateDelete())))
	mux.Handle("/books", (_mw.New(_mw.JSONResponse, _mw.Logger).Then(book.Create())))
	mux.Handle("/books/", (_mw.New(_mw.JSONResponse, _mw.Logger).Then(book.GetUpdateDelete())))
}
