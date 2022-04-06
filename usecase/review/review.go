package review

import (
	"log"
	"net/http"
	_bookRepository "plain-go/public-library/datastore/book"
	_userRepository "plain-go/public-library/datastore/user"
	_entity "plain-go/public-library/entity"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
	"strings"
	"time"
)

type ReviewUseCase struct {
	bookRepo _bookRepository.Book
	userRepo _userRepository.User
}

func New(book _bookRepository.Book, user _userRepository.User) *ReviewUseCase {
	return &ReviewUseCase{bookRepo: book, userRepo: user}
}

func (ruc ReviewUseCase) GetAllReviews() (res _model.GetAllReviewsResponse, code int, message string) {
	// calling repository
	reviews, err := ruc.bookRepo.GetAllReviews()

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, review := range reviews {
		// get reviewer detail
		user, err := ruc.userRepo.GetUserById(review.User.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// get book detail
		book, err := ruc.bookRepo.GetBookById(review.Book.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// formatting response
		user.Password = ""
		user.CreatedAt, _ = _helper.TimeFormatter(user.CreatedAt)
		user.UpdatedAt, _ = _helper.TimeFormatter(user.UpdatedAt)
		book.CreatedAt, _ = _helper.TimeFormatter(book.CreatedAt)
		book.UpdatedAt, _ = _helper.TimeFormatter(book.UpdatedAt)
		review.User = user
		review.Book = book
		res.Reviews = append(res.Reviews, review)
	}

	res.Reviews = reviews
	code, message = http.StatusOK, "success get all reviews"

	return
}

func (ruc ReviewUseCase) CreateReview(userId uint, bookId uint, req _model.CreateReviewRequest) (res _model.CreateReviewResponse, code int, message string) {
	// prepare input string
	content := strings.TrimSpace(req.Content)

	// check if required input is empty
	if content == "" {
		log.Println("empty input")
		code, message = http.StatusBadRequest, "empty input"
		return
	}

	// check if there is any forbidden character in required field
	if strings.Contains(strings.ReplaceAll(content, " ", ""), ";--") {
		log.Println("forbidden character")
		code, message = http.StatusBadRequest, "forbidden character"
		return
	}

	// check if star is out of range
	if req.Star < 1 || req.Star > 5 {
		log.Println("star out of range")
		code, message = http.StatusBadRequest, "star must be from 1 to 5"
		return
	}

	// check reviewer existence
	user, err := ruc.userRepo.GetUserById(userId)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if user.Name == "" {
		log.Println("user not found")
		code, message = http.StatusNotFound, "user not found"
		return
	}

	// check book existence
	book, err := ruc.bookRepo.GetBookById(bookId)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if book.Title == "" {
		log.Println("book not found")
		code, message = http.StatusNotFound, "book not found"
		return
	}

	// check if review already made
	reviews, err := ruc.bookRepo.GetAllReviewsByBookId(bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, review := range reviews {
		// member can only make 1 review per book
		if review.User.Id == userId {
			log.Println("review already exist")
			code, message = http.StatusConflict, "review already exist"
			return
		}
	}

	// prepare input to repository
	now := time.Now()
	newReview := _entity.SimplifiedReview{}
	newReview.User.Id = userId
	newReview.Book.Id = bookId
	newReview.Content = content
	newReview.Star = req.Star
	newReview.CreatedAt = now
	newReview.UpdatedAt = now

	// calling repository
	res.Review, err = ruc.bookRepo.CreateReview(newReview)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// formatting response
	res.Review.Book = book
	res.Review.Book.CreatedAt, _ = _helper.TimeFormatter(res.Review.Book.CreatedAt)
	res.Review.Book.UpdatedAt, _ = _helper.TimeFormatter(res.Review.Book.UpdatedAt)
	res.Review.User = user
	res.Review.User.CreatedAt, _ = _helper.TimeFormatter(res.Review.User.CreatedAt)
	res.Review.User.UpdatedAt, _ = _helper.TimeFormatter(res.Review.User.UpdatedAt)
	res.Review.CreatedAt, _ = _helper.TimeFormatter(res.Review.CreatedAt)
	res.Review.UpdatedAt, _ = _helper.TimeFormatter(res.Review.UpdatedAt)
	code, message = http.StatusCreated, "success create review"

	return
}

func (ruc ReviewUseCase) GetReviewByReviewId(bookId uint, reviewId uint) (res _model.GetReviewByIdResponse, code int, message string) {
	// calling repository
	review, err := ruc.bookRepo.GetReviewByReviewId(reviewId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// get reviewer
	review.User, err = ruc.userRepo.GetUserById(review.User.Id)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// get book
	review.Book, err = ruc.bookRepo.GetBookById(review.Book.Id)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if review.Book.Id != bookId {
		log.Println("this book has no such review")
		code, message = http.StatusNotFound, "this book has no such review"
		return
	}

	if review.Content == "" {
		log.Println("review not found")
		code, message = http.StatusNotFound, "review not found"
		return
	}

	// formatting response
	res.Review = review
	res.Review.CreatedAt, _ = _helper.TimeFormatter(res.Review.CreatedAt)
	res.Review.User.Password = ""
	res.Review.User.CreatedAt, _ = _helper.TimeFormatter(res.Review.User.CreatedAt)
	res.Review.User.UpdatedAt, _ = _helper.TimeFormatter(res.Review.User.UpdatedAt)
	res.Review.Book.CreatedAt, _ = _helper.TimeFormatter(res.Review.Book.CreatedAt)
	res.Review.Book.UpdatedAt, _ = _helper.TimeFormatter(res.Review.Book.UpdatedAt)
	code, message = http.StatusOK, "success get review by review id"

	return
}

func (ruc ReviewUseCase) GetAllReviewsByBookId(bookId uint) (res _model.GetAllReviewsByBookIdResponse, code int, message string) {
	// check book existence
	book, err := ruc.bookRepo.GetBookById(bookId)

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

	// calling repository
	reviews, err := ruc.bookRepo.GetAllReviewsByBookId(bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, review := range reviews {
		// get reviewer
		review.User, err = ruc.userRepo.GetUserById(review.User.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		review.User.Password = ""
		review.User.CreatedAt, _ = _helper.TimeFormatter(review.User.CreatedAt)
		review.User.UpdatedAt, _ = _helper.TimeFormatter(review.User.UpdatedAt)
		review.CreatedAt, _ = _helper.TimeFormatter(review.CreatedAt)
		review.UpdatedAt, _ = _helper.TimeFormatter(review.UpdatedAt)
	}

	res.Reviews = reviews
	res.Book = book
	code, message = http.StatusOK, "success get all reviews"

	return
}

func (ruc ReviewUseCase) UpdateReview(userId uint, bookId uint, reviewId uint, req _model.UpdateReviewRequest) (res _model.UpdateReviewResponse, code int, message string) {
	// check review existence
	review, err := ruc.bookRepo.GetReviewByReviewId(reviewId)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if review.Book.Id != bookId {
		log.Println("this book has no such review")
		code, message = http.StatusNotFound, "this book has no such review"
		return
	}

	// check if reviewer id not match
	if review.User.Id != userId {
		log.Println("forbidden")
		code, message = http.StatusForbidden, "forbidden"
		return
	}

	// prepare input string
	content := strings.TrimSpace(req.Content)

	flag := true

	if strings.Contains(strings.ReplaceAll(content, " ", ""), ";--") {
		log.Println("forbidden character")
		code, message = http.StatusBadRequest, "forbidden chacarter"
		return
	}

	if content != "" && content != review.Content {
		review.Content = content
		flag = false
	}

	// if req.Star = 0, it means no change
	if req.Star > 5 {
		log.Println("star out of range")
		code, message = http.StatusBadRequest, "star must be from 1 to 5"
		return
	}

	if req.Star != 0 && req.Star != review.Star {
		review.Star = req.Star
		flag = false
	}

	if flag {
		log.Println("no update was performed")
		code, message = http.StatusBadRequest, "no update was performed"
		return
	}

	// calling repository
	res.Review, err = ruc.bookRepo.UpdateReview(review)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// get reviewer
	res.Review.User, err = ruc.userRepo.GetUserById(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// get book
	res.Review.Book, err = ruc.bookRepo.GetBookById(res.Review.Book.Id)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// formatting response
	res.Review.User.Password = ""
	res.Review.User.CreatedAt, _ = _helper.TimeFormatter(res.Review.User.CreatedAt)
	res.Review.User.UpdatedAt, _ = _helper.TimeFormatter(res.Review.User.UpdatedAt)
	res.Review.Book.CreatedAt, _ = _helper.TimeFormatter(res.Review.Book.CreatedAt)
	res.Review.Book.UpdatedAt, _ = _helper.TimeFormatter(res.Review.Book.UpdatedAt)
	res.Review.CreatedAt, _ = _helper.TimeFormatter(res.Review.CreatedAt)
	res.Review.UpdatedAt, _ = _helper.TimeFormatter(res.Review.UpdatedAt)
	code, message = http.StatusOK, "success update review"

	return
}

func (ruc ReviewUseCase) UpdateStatus(bookId uint, reviewId uint, req _model.UpdateReviewRequest) (code int, message string) {
	// check review existence
	review, err := ruc.bookRepo.GetReviewByReviewId(reviewId)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if review.Book.Id != bookId {
		log.Println("this book has no such review")
		code, message = http.StatusNotFound, "this book has no such review"
		return
	}

	flag := uint(0)
	if req.IsRead {
		flag = 1
	}

	if review.Flag == flag {
		log.Println("no update was performed")
		code, message = http.StatusBadRequest, "no update was performed"
		return
	}

	// calling repository
	if err = ruc.bookRepo.UpdateReviewStatus(flag, reviewId); err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	code, message = http.StatusOK, "success update review status"

	return
}

func (ruc ReviewUseCase) DeleteReview(userId uint, bookId uint, reviewId uint) (code int, message string) {
	// check review existence
	review, err := ruc.bookRepo.GetReviewByReviewId(reviewId)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if review.Book.Id != bookId {
		log.Println("this book has no such review")
		code, message = http.StatusNotFound, "this book has no such review"
		return
	}

	// check if reviewer id not match
	if review.User.Id != userId {
		log.Println("forbidden")
		code, message = http.StatusForbidden, "forbidden"
		return
	}

	if err = ruc.bookRepo.DeleteReview(reviewId); err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	code, message = http.StatusOK, "success delete review"

	return
}
