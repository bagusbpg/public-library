package request

import (
	"net/http"
	_bookRepository "plain-go/public-library/datastore/book"
	_requestRepository "plain-go/public-library/datastore/request"
	_userRepository "plain-go/public-library/datastore/user"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
)

type RequestUseCase struct {
	bookRepo    _bookRepository.Book
	userRepo    _userRepository.User
	requestRepo _requestRepository.Request
}

func New(book _bookRepository.Book, user _userRepository.User, request _requestRepository.Request) *RequestUseCase {
	return &RequestUseCase{bookRepo: book, userRepo: user, requestRepo: request}
}

func (ruc RequestUseCase) GetAllRequests() (res _model.GetAllRequestResponse, code int, message string) {
	// calling repository
	requests, err := ruc.requestRepo.GetAllRequests()

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, request := range requests {
		// get reviewer detail
		request.User, err = ruc.userRepo.GetUserById(request.User.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// omit deleted reviewer
		if request.User.Name == "" {
			continue
		}

		// get book detail
		request.BookItem.Book, err = ruc.bookRepo.GetBookByItemId(request.BookItem.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// omit deleted book
		if request.BookItem.Book.Title == "" {
			continue
		}

		averageStar, err := ruc.bookRepo.CountStarsByBookId(request.BookItem.Book.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// formatting response
		request.User.CreatedAt, _ = _helper.TimeFormatter(request.User.CreatedAt)
		request.User.UpdatedAt, _ = _helper.TimeFormatter(request.User.UpdatedAt)
		request.BookItem.Book.CreatedAt, _ = _helper.TimeFormatter(request.BookItem.Book.CreatedAt)
		request.BookItem.Book.UpdatedAt, _ = _helper.TimeFormatter(request.BookItem.Book.UpdatedAt)
		request.BookItem.Book.AverageStar = _helper.NilHandler(averageStar)
		request.CreatedAt, _ = _helper.TimeFormatter(request.CreatedAt)
		request.UpdatedAt, _ = _helper.TimeFormatter(request.UpdatedAt)

		res.Requests = append(res.Requests, request)
	}

	code, message = http.StatusOK, "success get all requests"

	return
}

func (ruc RequestUseCase) GetAllRequestsByUserId(userId uint) (res _model.GetAllRequestResponse, code int, message string) {

	return
}
