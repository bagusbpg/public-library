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
	"time"

	// _entity "plain-go/public-library/entity"
	_model "plain-go/public-library/model"
)

type WishUseCase struct {
	bookRepo _bookRepository.Book
	userRepo _userRepository.User
}

func New(book _bookRepository.Book, user _userRepository.User) *WishUseCase {
	return &WishUseCase{bookRepo: book, userRepo: user}
}

func (wuc WishUseCase) AddBookToWishlist(userId uint, req _model.AddBookToWishlistRequest) (res _model.AddBookToWishlistResponse, code int, message string) {
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

	// check if book already purchased
	book, err := wuc.bookRepo.GetBookByTitle(req.Title)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if book.Title != "" {
		log.Println("book already exist")
		code, message = http.StatusConflict, "book already exist"
		return
	}

	// check if book already in wishlist
	wishes, err := wuc.bookRepo.GetWishesByUserId(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, wish := range wishes {
		if wish.Title == title {
			code, message = http.StatusConflict, "book already in wishlist"
			return
		}
	}

	// prepare input to repository
	newWish := _entity.Wish{}
	newWish.Title = title
	newWish.Category = category
	newWish.CreatedAt = time.Now()

	// calling repository
	res.Wish, err = wuc.bookRepo.AddBookToWishlist(userId, newWish)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// create author
	for _, _author := range req.Author {
		// calling repository
		author, err := wuc.bookRepo.GetAuthorByName(_author.Name)

		// detect failure in repository
		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// if author does not exist, create new
		if author.Name == "" {
			author.Name = _author.Name

			// calling repository
			author, err = wuc.bookRepo.CreateNewAuthor(author)

			// detect failure in repository
			if err != nil {
				code, message = http.StatusInternalServerError, "internal server error"
				return
			}
		}

		// if author exist or after author created, create book author junction
		if err = wuc.bookRepo.CreateWishAuthorJunction(newWish, author); err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		res.Wish.Author = append(res.Wish.Author, author)
	}

	res.Wish.CreatedAt, _ = _helper.TimeFormatter(res.Wish.CreatedAt)
	code, message = http.StatusCreated, "success add book to wishlist"

	return
}

func (wuc WishUseCase) RemoveBookFromWishlist(userId uint, wishId uint) (code int, message string) {
	// check user existence
	user, err := wuc.userRepo.GetUserById(userId)

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

	// check wish existence
	wish, err := wuc.bookRepo.GetWishById(userId, wishId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if wish.Title == "" {
		log.Println("wish not found")
		code, message = http.StatusNotFound, "wish not found"
		return
	}

	// calling repository
	if err = wuc.bookRepo.RemoveBookFromWishlist(wishId); err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	code, message = http.StatusOK, "success remove book from favorites"

	return
}

func (wuc WishUseCase) GetAllWishes() (res []_model.GetWishesByUserIdResponse, code int, message string) {
	// get all users
	users, err := wuc.userRepo.GetAllUsers()

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, user := range users {
		// get wishes of each user
		wishes, err := wuc.bookRepo.GetWishesByUserId(user.Id)

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

		}

		if len(wishes) != 0 {
			// formatting response
			user.Password = ""
			user.CreatedAt, _ = _helper.TimeFormatter(user.CreatedAt)
			user.UpdatedAt, _ = _helper.TimeFormatter(user.UpdatedAt)
			res = append(res, _model.GetWishesByUserIdResponse{User: user, Wishes: wishes})
		}
	}

	code, message = http.StatusOK, "success get all wishes"

	return
}

func (wuc WishUseCase) GetAllWishesByUserId(userId uint) (res _model.GetWishesByUserIdResponse, code int, message string) {
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
	wishes, err := wuc.bookRepo.GetWishesByUserId(userId)

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
