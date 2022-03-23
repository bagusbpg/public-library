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
	if req.Star > 5 {
		log.Println("star out of range")
		code, message = http.StatusBadRequest, "star cannot be greater than 5"
		return
	}

	// check user existence
	user, err := ruc.userRepo.GetUserById(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// check if user does not exist
	if user.Name == "" {
		log.Println("user not found")
		code, message = http.StatusNotFound, "user not found"
		return
	}

	// check book existence
	book, err := ruc.bookRepo.GetBookById(bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// check if book does not exist
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
	newReview := _entity.Review{}
	newReview.Content = content
	newReview.Star = req.Star
	newReview.CreatedAt = now
	newReview.UpdatedAt = now

	// calling repository
	res.Review, err = ruc.bookRepo.CreateReview(userId, bookId, newReview)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// formatting response
	res.Book = book
	res.Book.CreatedAt, _ = _helper.TimeFormatter(res.Book.CreatedAt)
	res.Book.UpdatedAt, _ = _helper.TimeFormatter(res.Book.UpdatedAt)
	res.Review.User = user
	res.Review.User.CreatedAt, _ = _helper.TimeFormatter(res.Review.User.CreatedAt)
	res.Review.User.UpdatedAt, _ = _helper.TimeFormatter(res.Review.User.UpdatedAt)
	res.Review.CreatedAt, _ = _helper.TimeFormatter(res.Review.CreatedAt)
	res.Review.UpdatedAt, _ = _helper.TimeFormatter(res.Review.UpdatedAt)
	code, message = http.StatusCreated, "success create review"

	return
}

func (ruc ReviewUseCase) UpdateReview() {

}
