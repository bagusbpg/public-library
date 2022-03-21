package wish

import (
	// "log"
	// "net/http"
	"log"
	"net/http"
	_bookRepository "plain-go/public-library/datastore/book"
	_userRepository "plain-go/public-library/datastore/user"
	_entity "plain-go/public-library/entity"
	_helper "plain-go/public-library/helper"
	"strings"

	// _entity "plain-go/public-library/entity"
	_model "plain-go/public-library/model"
)

type WishUseCase struct {
	userRepo _userRepository.User
	bookRepo _bookRepository.Book
}

func New(user _userRepository.User, book _bookRepository.Book) *WishUseCase {
	return &WishUseCase{userRepo: user, bookRepo: book}
}

func (wuc WishUseCase) AddBookToWishlist(userId uint, req _model.AddBookToWishlistRequest) (code int, message string) {
	// prepare input string
	title := strings.Title(strings.TrimSpace(req.Title))
	category := strings.TrimSpace(req.Category)

	check := []string{title, category}

	for _, s := range check {
		// check if required input is empty
		if s == "" {
			log.Println("empty input")
			code, message = http.StatusBadRequest, "empty input"
			return
		}

		// check if there is any forbidden character in required field
		if strings.Contains(strings.ReplaceAll(s, " ", ""), ";--") {
			log.Println("forbidden character")
			code, message = http.StatusBadRequest, "forbidden character"
			return
		}
	}

	// check authors
	if len(req.Author) == 0 {
		log.Println("empty input")
		code, message = http.StatusBadRequest, "empty input"
		return
	}

	for _, _author := range req.Author {
		_author.Name = strings.TrimSpace(_author.Name)

		// check if required input is empty
		if _author.Name == "" {
			log.Println("empty input")
			code, message = http.StatusBadRequest, "empty input"
			return
		}

		// check if there is any forbidden character in required field
		if strings.Contains(strings.ReplaceAll(_author.Name, " ", ""), ";--") {
			log.Println("forbidden character")
			code, message = http.StatusBadRequest, "forbidden character"
			return
		}
	}

	// check user existence
	user, err := wuc.userRepo.GetUserById(userId)

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

	// check if book already exist
	book, err := wuc.bookRepo.GetBookByTitle(req.Title)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// check if book is already exist (purchased by library)
	if book.Title != "" {
		log.Println("book already exist")
		code, message = http.StatusConflict, "book already exist"
		return
	}

	// check if book is already entered in wishlist

	return
}

func (wuc WishUseCase) GetAllWishes(userId uint) (res _model.GetAllWishes, code int, message string) {
	// check user existence
	user, err := wuc.userRepo.GetUserById(userId)

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

	// calling repository
	wishes, err := wuc.bookRepo.GetAllWishes(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, wish := range wishes {
		// getting each book's authors
		authors, err := wuc.bookRepo.GetWishAuthors(wish.Id)

		// detect failure in repository
		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		wish.Author = authors
		wish.CreatedAt, _ = _helper.TimeFormatter(wish.CreatedAt)

		res.Wishes = append(res.Wishes, _entity.Wish{Id: wish.Id, Title: wish.Title, Author: wish.Author, Category: wish.Category, CreatedAt: wish.CreatedAt})
	}

	res.User = user
	res.User.Password = ""
	res.User.CreatedAt, _ = _helper.TimeFormatter(res.User.CreatedAt)
	res.User.UpdatedAt, _ = _helper.TimeFormatter(res.User.UpdatedAt)
	code, message = http.StatusOK, "success get all wishes"

	return
}
