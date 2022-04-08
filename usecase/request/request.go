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
		request.User.Password = ""
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
		request.StartAt, _ = _helper.TimeFormatter(request.StartAt)
		request.FinishAt, _ = _helper.TimeFormatter(request.FinishAt)
		request.ReturnAt, _ = _helper.TimeFormatter(request.ReturnAt)
		request.CancelAt, _ = _helper.TimeFormatter(request.CancelAt)
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
		request.StartAt, _ = _helper.TimeFormatter(request.StartAt)
		request.FinishAt, _ = _helper.TimeFormatter(request.FinishAt)
		request.ReturnAt, _ = _helper.TimeFormatter(request.ReturnAt)
		request.CancelAt, _ = _helper.TimeFormatter(request.CancelAt)
		request.UpdatedAt, _ = _helper.TimeFormatter(request.UpdatedAt)

		res.Requests = append(res.Requests, request)
	}

	code, message = http.StatusOK, "success get all requests"

	return
}

func (ruc RequestUseCase) GetRequestById(userId uint, requestId uint) (res _model.GetRequestByIdResponse, code int, message string) {
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

	// check if requester does not match
	if request.User.Id != userId {
		log.Println("forbidden")
		code, message = http.StatusForbidden, "forbidden"
		return
	}

	// get requester
	request.User, err = ruc.userRepo.GetUserById(userId)

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
	request.User.Password = ""
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
	request.StartAt, _ = _helper.TimeFormatter(request.StartAt)
	request.FinishAt, _ = _helper.TimeFormatter(request.FinishAt)
	request.ReturnAt, _ = _helper.TimeFormatter(request.ReturnAt)
	request.CancelAt, _ = _helper.TimeFormatter(request.CancelAt)
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
		newRequest.Status.Id = 1
		newRequest.Status.Description = "waiting in queue"
	} else {
		newRequest.BookItem.Id = int(bookItemId)
		newRequest.Status.Id = 2
		newRequest.Status.Description = "book is being prepared"
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
	user.Password = ""
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

func (ruc RequestUseCase) UpdateRequest(userId uint, requestId uint, role string, req _model.UpdateRequestRequest) (res _model.UpdateRequestResponse, code int, message string) {
	// check request existence
	request, err := ruc.requestRepo.GetRequestById(requestId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if request.Status.Id == 0 {
		code, message = http.StatusNotFound, "request not found"
		return
	}

	if request.CancelAt != nil {
		code, message = http.StatusBadRequest, "request already cancelled"
		return
	}

	// global timestamp
	now := time.Now()
	request.UpdatedAt = now

	switch role {
	case "Member":
		// check if requester does not match
		if request.User.Id != userId {
			log.Println("forbidden")
			code, message = http.StatusForbidden, "forbidden"
			return
		}

		switch req.ActionCode {
		case 11: // cancel request
			// check if cancelling request is possible
			if request.Status.Id > 2 {
				log.Println("cannot cancel request at this time")
				code, message = http.StatusBadRequest, "cannot cancel request at this time"
				return
			}

			// prepare input to repository
			request.Status.Id = 3 // "request is cancelled"
			request.CancelAt = now
		case 12: // extend request
			// check if extending request is possible
			if request.Extended != 0 || request.Status.Id != 5 {
				log.Println("cannot extend request at this time")
				code, message = http.StatusBadRequest, "cannot extend request at this time"
				return
			}

			// prepare input to repository
			request.Status.Id = 6 // "request is extended"
			request.Extended = 1
			request.FinishAt = now.AddDate(0, 0, 7)
		}
	case "Librarian":
		switch req.ActionCode {
		case 21: // notify for book pick up
			// check if request is not cancelled and notifying is possible
			if request.Status.Id != 2 {
				log.Println("cannot notify pick up at this time")
				code, message = http.StatusBadRequest, "cannot notify pick up at this time"
				return
			}

			// prepare input to repository
			request.Status.Id = 4 // "book is ready for pick up"
		case 22: // borrow period started
			// check if request is not cancelled and borrowing is possible
			if request.Status.Id != 4 {
				log.Println("cannot hand over book at this time")
				code, message = http.StatusBadRequest, "cannot hand over book at this time"
				return
			}

			// prepare input to repository
			request.Status.Id = 5 // "book is borrowed"
			request.FinishAt = now.AddDate(0, 0, 7)
		case 23: // normal return (no penalty)
			// check if normal return is possible
			if request.Status.Id < 5 || request.Status.Id > 6 {
				log.Println("return is not possible at this time")
				code, message = http.StatusBadRequest, "return is not possible at this time"
				return
			}

			if finishAt := request.FinishAt.(time.Time); finishAt.After(now) {
				log.Println("penalty must be paid first")
				code, message = http.StatusBadRequest, "penalty must be paid first"
				return
			}

			// prepare input to repository
			request.Status.Id = 8
			request.ReturnAt = now
		case 24: // late return (with penalty)
			if request.Status.Id != 7 {
				log.Println("penalty is not payable at this time")
				code, message = http.StatusBadRequest, "penalty is not payable at this time"
				return
			}

			// prepare input to repository
			request.Status.Id = 9
			request.ReturnAt = now
		}
	default:
		log.Println("role not assigned")
		code, message = http.StatusBadRequest, "role not assigned"
	}

	// calling repository
	res.Request, err = ruc.requestRepo.Update(request)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// check for late return
	if res.Request.Status.Id == 5 || res.Request.Status.Id == 6 {
		borrowDuration := time.Until(res.Request.FinishAt.(time.Time))
		afterFunctTimer := time.AfterFunc((borrowDuration), func() {
			futureReq, _ := ruc.requestRepo.GetRequestById(res.Request.Id)

			if futureReq.Status.Id == 5 || futureReq.Status.Id == 6 {
				futureReq.Status.Id = 7
				futureReq.UpdatedAt = futureReq.FinishAt.(time.Time)
				ruc.requestRepo.Update(futureReq)
			}
		})

		defer afterFunctTimer.Stop()

		timer := time.NewTimer(borrowDuration)
		<-timer.C
	}

	// formatting response
	res.Request.User, _ = ruc.userRepo.GetUserById(res.Request.User.Id)
	res.Request.User.Password = ""
	res.Request.User.CreatedAt, _ = _helper.TimeFormatter(res.Request.User.CreatedAt)
	res.Request.User.UpdatedAt, _ = _helper.TimeFormatter(res.Request.UpdatedAt)
	res.Request.BookItem.Book, _ = ruc.bookRepo.GetBookByItemId(uint(res.Request.BookItem.Id))
	res.Request.BookItem.Book.CreatedAt, _ = _helper.TimeFormatter(res.Request.BookItem.Book.CreatedAt)
	res.Request.BookItem.Book.UpdatedAt, _ = _helper.TimeFormatter(res.Request.BookItem.Book.UpdatedAt)
	res.Request.BookItem.Book.Quantity, _ = ruc.bookRepo.CountBookById(res.Request.BookItem.Book.Id)
	res.Request.BookItem.Book.Author, _ = ruc.bookRepo.GetBookAuthors(res.Request.BookItem.Book.Id)
	res.Request.BookItem.Book.FavoriteCount, _ = ruc.bookRepo.CountFavoritesByBookId(res.Request.BookItem.Book.Id)
	averageStar, _ := ruc.bookRepo.CountStarsByBookId(res.Request.BookItem.Book.Id)
	res.Request.BookItem.Book.AverageStar = _helper.NilHandler(averageStar)
	res.Request.CreatedAt, _ = _helper.TimeFormatter(res.Request.CreatedAt)
	res.Request.UpdatedAt, _ = _helper.TimeFormatter(res.Request.UpdatedAt)
	if res.Request.CancelAt != nil {
		res.Request.CancelAt, _ = _helper.TimeFormatter(res.Request.CancelAt)
	}
	code, message = http.StatusOK, "success update request"

	return
}
