package favorite

import (
	"log"
	"net/http"
	_bookRepository "plain-go/public-library/datastore/book"
	_userRepository "plain-go/public-library/datastore/user"
	_entity "plain-go/public-library/entity"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
)

type FavoriteUseCase struct {
	bookRepo _bookRepository.Book
	userRepo _userRepository.User
}

func New(book _bookRepository.Book, user _userRepository.User) *FavoriteUseCase {
	return &FavoriteUseCase{bookRepo: book, userRepo: user}
}

func (fuc FavoriteUseCase) AddBookToFavorite(userId uint, bookId uint) (code int, message string) {
	// check user existence
	user, err := fuc.userRepo.GetUserById(userId)

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
	book, err := fuc.bookRepo.GetBookById(bookId)

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

	// check if book already in favorites
	favs, err := fuc.bookRepo.GetAllFavorites(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	favMap := map[uint]interface{}{}

	for _, item := range favs {
		favMap[item.Book.Id] = nil
	}

	if _, exist := favMap[bookId]; exist {
		code, message = http.StatusConflict, "book already in favorites"
		return
	}

	// calling repository
	if err = fuc.bookRepo.AddBookToFavorite(userId, bookId); err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	code, message = http.StatusCreated, "success add book to favorites"

	return
}

func (fuc FavoriteUseCase) RemoveBookFromFavorite(userId uint, bookId uint) (code int, message string) {
	// check user existence
	user, err := fuc.userRepo.GetUserById(userId)

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
	book, err := fuc.bookRepo.GetBookById(bookId)

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

	// check if book already in favorites
	favs, err := fuc.bookRepo.GetAllFavorites(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	favMap := map[uint]interface{}{}

	for _, item := range favs {
		favMap[item.Book.Id] = nil
	}

	if _, exist := favMap[bookId]; !exist {
		code, message = http.StatusNotFound, "book not in favorites"
		return
	}

	// calling repository
	if err = fuc.bookRepo.RemoveBookFromFavorite(userId, bookId); err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	code, message = http.StatusOK, "success remove book from favorites"

	return
}

func (fuc FavoriteUseCase) GetAllFavorites(userId uint) (res _model.GetAllFavoritesResponse, code int, message string) {
	// check user existence
	user, err := fuc.userRepo.GetUserById(userId)

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
	favs, err := fuc.bookRepo.GetAllFavorites(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, fav := range favs {
		// getting each book
		book, err := fuc.bookRepo.GetBookById(fav.Book.Id)

		// detect failure in repository
		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// getting each book's authors
		authors, err := fuc.bookRepo.GetBookAuthors(book.Id)

		// detect failure in repository
		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		book.Author = authors
		fav.CreatedAt, _ = _helper.TimeFormatter(fav.CreatedAt)

		res.Favorites = append(res.Favorites, _entity.Favorite{Id: fav.Id, Book: book, CreatedAt: fav.CreatedAt})
	}

	code, message = http.StatusOK, "success get all favorites"

	return
}
