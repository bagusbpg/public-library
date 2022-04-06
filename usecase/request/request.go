package request

import (
	"log"
	"net/http"
	_bookRepository "plain-go/public-library/datastore/book"
	_requestRepository "plain-go/public-library/datastore/request"
	_userRepository "plain-go/public-library/datastore/user"
	_entity "plain-go/public-library/entity"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
	"time"
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
		request.BookItem.Book, err = ruc.bookRepo.GetBookByItemId(uint(request.BookItem.Id))

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// omit deleted book
		if request.BookItem.Book.Title == "" {
			continue
		}

		// formatting response
		request.User.CreatedAt, _ = _helper.TimeFormatter(request.User.CreatedAt)
		request.User.UpdatedAt, _ = _helper.TimeFormatter(request.User.UpdatedAt)
		request.BookItem.Book.CreatedAt, _ = _helper.TimeFormatter(request.BookItem.Book.CreatedAt)
		request.BookItem.Book.UpdatedAt, _ = _helper.TimeFormatter(request.BookItem.Book.UpdatedAt)
		request.BookItem.Book.Quantity, _ = ruc.bookRepo.CountBookById(request.BookItem.Book.Id)
		request.BookItem.Book.Author, _ = ruc.bookRepo.GetBookAuthors(request.BookItem.Book.Id)
		request.BookItem.Book.FavoriteCount, _ = ruc.bookRepo.CountFavoritesByBookId(request.BookItem.Book.Id)
		averageStar, _ := ruc.bookRepo.CountStarsByBookId(request.BookItem.Book.Id)
		request.BookItem.Book.AverageStar = _helper.NilHandler(averageStar)
		request.CreatedAt, _ = _helper.TimeFormatter(request.CreatedAt)
		request.UpdatedAt, _ = _helper.TimeFormatter(request.UpdatedAt)

		res.Requests = append(res.Requests, request)
	}

	code, message = http.StatusOK, "success get all requests"

	return
}

func (ruc RequestUseCase) GetAllRequestsByUserId(userId uint) (res _model.GetAllRequestByUserIdResponse, code int, message string) {
	// check requester existence
	user, err := ruc.userRepo.GetUserById(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if user.Name == "" {
		log.Println("user not found")
		code, message = http.StatusNotFound, "user not found"
		return
	}

	// calling repository
	requests, err := ruc.requestRepo.GetAllRequestsByUserId(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, request := range requests {
		// formatting response
		request.BookItem.Book.CreatedAt, _ = _helper.TimeFormatter(request.BookItem.Book.CreatedAt)
		request.BookItem.Book.UpdatedAt, _ = _helper.TimeFormatter(request.BookItem.Book.UpdatedAt)
		request.BookItem.Book.Quantity, _ = ruc.bookRepo.CountBookById(request.BookItem.Book.Id)
		request.BookItem.Book.Author, _ = ruc.bookRepo.GetBookAuthors(request.BookItem.Book.Id)
		request.BookItem.Book.FavoriteCount, _ = ruc.bookRepo.CountFavoritesByBookId(request.BookItem.Book.Id)
		averageStar, _ := ruc.bookRepo.CountStarsByBookId(request.BookItem.Book.Id)
		request.BookItem.Book.AverageStar = _helper.NilHandler(averageStar)
		request.CreatedAt, _ = _helper.TimeFormatter(request.CreatedAt)
		request.UpdatedAt, _ = _helper.TimeFormatter(request.UpdatedAt)

		res.Requests = append(res.Requests, request)
	}

	code, message = http.StatusOK, "success get all requests"

	return
}

func (ruc RequestUseCase) GetRequestById(requestId uint) (res _model.GetRequestByIdResponse, code int, message string) {
	// calling repository
	request, err := ruc.requestRepo.GetRequestById(requestId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if request.Status.Id == 0 {
		log.Println("request not found")
		code, message = http.StatusNotFound, "request not found"
		return
	}

	// get requester
	request.User, err = ruc.userRepo.GetUserById(request.User.Id)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// get book
	request.BookItem.Book, err = ruc.bookRepo.GetBookByItemId(uint(request.BookItem.Id))

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// formatting response
	request.User.CreatedAt, _ = _helper.TimeFormatter(request.User.CreatedAt)
	request.User.UpdatedAt, _ = _helper.TimeFormatter(request.User.UpdatedAt)
	request.BookItem.Book.CreatedAt, _ = _helper.TimeFormatter(request.BookItem.Book.CreatedAt)
	request.BookItem.Book.UpdatedAt, _ = _helper.TimeFormatter(request.BookItem.Book.UpdatedAt)
	request.BookItem.Book.Quantity, _ = ruc.bookRepo.CountBookById(request.BookItem.Book.Id)
	request.BookItem.Book.Author, _ = ruc.bookRepo.GetBookAuthors(request.BookItem.Book.Id)
	request.BookItem.Book.FavoriteCount, _ = ruc.bookRepo.CountFavoritesByBookId(request.BookItem.Book.Id)
	averageStar, _ := ruc.bookRepo.CountStarsByBookId(request.BookItem.Book.Id)
	request.BookItem.Book.AverageStar = _helper.NilHandler(averageStar)
	request.CreatedAt, _ = _helper.TimeFormatter(request.CreatedAt)
	request.UpdatedAt, _ = _helper.TimeFormatter(request.UpdatedAt)

	return
}

func (ruc RequestUseCase) CreateRequest(userId uint, req _model.CreateRequestRequest) (res _model.CreateRequestResponse, code int, message string) {
	// check requester existence
	user, err := ruc.userRepo.GetUserById(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if user.Name == "" {
		log.Println("user not found")
		code, message = http.StatusNotFound, "user not found"
		return
	}

	// check request limit
	count, err := ruc.requestRepo.CountActiveRequestByUserId(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if count >= 2 {
		code, message = http.StatusForbidden, "requests reached maximum limit"
		return
	}

	// check book existence
	book, err := ruc.bookRepo.GetBookById(req.BookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if book.Title == "" {
		log.Println("book not found")
		code, message = http.StatusNotFound, "book not found"
		return
	}

	// check book availability
	bookItemId, err := ruc.bookRepo.GetAvailableBookByBookId(req.BookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// prepare input to repository
	now := time.Now()
	newRequest := _entity.Request{}
	newRequest.User.Id = userId

	if bookItemId == 0 {
		newRequest.BookItem.Id = -1
		newRequest.Status.Description = "waiting in queue"
	} else {
		newRequest.BookItem.Id = int(bookItemId)
		newRequest.Status.Description = "book is being prepared"
	}

	newRequest.Status.Id, err = ruc.requestRepo.GetRequestStatusId(newRequest.Status.Description)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if newRequest.Status.Id == 0 {
		code, message = http.StatusNotFound, "status not found"
		return
	}

	newRequest.CreatedAt = now
	newRequest.UpdatedAt = now

	// calling repository
	res.Request, err = ruc.requestRepo.CreateNewRequest(newRequest)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// formatting response
	user.CreatedAt, _ = _helper.TimeFormatter(user.CreatedAt)
	user.UpdatedAt, _ = _helper.TimeFormatter(user.UpdatedAt)
	res.Request.User = user

	book.CreatedAt, _ = _helper.TimeFormatter(book.CreatedAt)
	book.UpdatedAt, _ = _helper.TimeFormatter(book.UpdatedAt)
	averageStar, _ := ruc.bookRepo.CountStarsByBookId(req.BookId)
	book.AverageStar = _helper.NilHandler(averageStar)
	book.FavoriteCount, _ = ruc.bookRepo.CountFavoritesByBookId(req.BookId)
	book.Author, _ = ruc.bookRepo.GetBookAuthors(req.BookId)
	book.Quantity, _ = ruc.bookRepo.CountBookById(req.BookId)
	res.Request.BookItem.Book = book

	res.Request.CreatedAt, _ = _helper.TimeFormatter(res.Request.CreatedAt)
	res.Request.UpdatedAt, _ = _helper.TimeFormatter(res.Request.UpdatedAt)

	code, message = http.StatusCreated, "success create new request"

	return
}

func (ruc RequestUseCase) UpdateRequest() {

}
