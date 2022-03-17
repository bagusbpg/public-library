package router

import (
	"net/http"
	_mw "plain-go/public-library/app/middleware"
	_book "plain-go/public-library/controller/book"
	_favorite "plain-go/public-library/controller/favorite"
	_user "plain-go/public-library/controller/user"
)

func Router(
	mux *http.ServeMux,
	user *_user.UserController,
	book *_book.BookController,
	favorite *_favorite.FavoriteController,
) {
	mux.Handle("/public/signup", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.POST, _mw.JSONRequest).Then(user.SignUp()))
	// mux.Handle("/public/stats", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.GET).Then(book.Stats()))
	mux.Handle("/public/books", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.GET).Then(book.GetAll()))
	mux.Handle("/public/books/", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.GET, _mw.ValidateId).Then(book.Get()))
	// mux.Handle("/public/reviews/", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.GET).Then(review.GetAll()))

	mux.Handle("/login", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.POST, _mw.JSONRequest).Then(user.Login()))

	mux.Handle("/member/profiles", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.GET, _mw.MemberAndLibrarianAuthorization).Then(user.Get()))
	mux.Handle("/member/profiles/", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.PUT, _mw.JSONRequest, _mw.ValidateId, _mw.MemberOnlydAuthorization).Then(user.Update()))
	// GET		member/requests -> get all my request (based on token); use query params to see history of past requests
	// GET		member/requests/ -> get my specific request by request id (verified by token)
	// POST		member/reviews -> create reviews (verified by token)
	// PUT 		member/reviews/1 -> update review by review id (verified by token)
	// POST		member/favorites -> create favorite book (based on token)
	mux.Handle("/member/favorites/", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.GET, _mw.ValidateId, _mw.MemberAndLibrarianAuthorization).Then(favorite.GetAllFavorites()))
	// OK GET	member/favorites/1 -> get all my favorite books by user id
	// POST		member/wishlist -> create wishlist
	// GET		member/wishlist/1 -> get all my wishlist by user id
	mux.Handle("/member/delete/profiles/", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.DELETE, _mw.ValidateId, _mw.MemberAndLibrarianAuthorization).Then(user.Delete()))
	// OK DELETE	member/delete/profiles/1 -> delete profile by user id (verified by token)
	// DELETE	member/delete/requests/1 -> cancel borrow request by request id (verified by token)
	// DELETE	member/delete/reviews/1 -> delete my review based on review id (verified by token)
	// DELETE	member/delete/favorites/1 -> delete my favorite based on book id (verified by token)

	mux.Handle("/librarian/books", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.POST, _mw.JSONRequest, _mw.LibrarianAuthorization).Then(book.Create()))
	mux.Handle("/librarian/books/", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.PUT, _mw.JSONRequest, _mw.LibrarianAuthorization).Then(book.Update()))
	mux.Handle("/librarian/users", _mw.New(_mw.JSONResponse, _mw.Logger, _mw.GET, _mw.LibrarianAuthorization).Then(user.GetAll()))
	// GET		librarian/requests -> get all requests
	// GET		librarian/requests/1 -> get request of specific user
	// GET		librarian/requests/1/1 -> get specific request of specific user
}
