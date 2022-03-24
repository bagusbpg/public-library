package favorite

import (
	"log"
	"net/http"
	_bookRepository "plain-go/public-library/datastore/book"
	_userRepository "plain-go/public-library/datastore/user"
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

func (fuc FavoriteUseCase) AddBookToFavorite(userId uint, bookId uint) (res _model.AddBookToFavoriteResponse, code int, message string) {
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
	favs, err := fuc.bookRepo.GetAllFavoritesByUserId(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, fav := range favs {
		if fav.Book.Id == bookId {
			code, message = http.StatusConflict, "book already in favorites"
			return
		}
	}

	// calling repository
	res.Favorite, err = fuc.bookRepo.AddBookToFavorite(userId, bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	res.Favorite.Book = book
	res.Favorite.CreatedAt, _ = _helper.TimeFormatter(res.Favorite.CreatedAt)
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

	if book.Title == "" {
		log.Println("book not found")
		code, message = http.StatusNotFound, "book not found"
		return
	}

	// check if book already in favorites
	favs, err := fuc.bookRepo.GetAllFavoritesByUserId(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	flag := true

	for i := 0; i < len(favs) && flag; i++ {
		if favs[i].Book.Id == bookId {
			flag = false
		}
	}

	if flag {
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

func (fuc FavoriteUseCase) GetAllFavoritesByUserId(userId uint) (res _model.GetAllFavoritesResponse, code int, message string) {
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
	favs, err := fuc.bookRepo.GetAllFavoritesByUserId(userId)

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
		fav.Book = book
		fav.CreatedAt, _ = _helper.TimeFormatter(fav.CreatedAt)

		res.Favorites = append(res.Favorites, fav)
	}

	res.User = user
	res.User.Password = ""
	res.User.CreatedAt, _ = _helper.TimeFormatter(res.User.CreatedAt)
	res.User.UpdatedAt, _ = _helper.TimeFormatter(res.User.UpdatedAt)
	code, message = http.StatusOK, "success get all favorites"

	return
}
