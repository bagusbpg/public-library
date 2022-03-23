package book

import (
	_entity "plain-go/public-library/entity"
)

type Book interface {
	GetBookByTitle(title string) (book _entity.Book, err error)
	GetAuthorByName(name string) (author _entity.Author, err error)
	GetAllAuthors() (authors []_entity.Author, err error)
	CreateNewBook(newBook _entity.Book) (book _entity.Book, err error)
	CreateNewAuthor(newAuthor _entity.Author) (author _entity.Author, err error)
	CreateBookAuthorJunction(book _entity.Book, author _entity.Author) (err error)
	CreateBookItem(book _entity.Book) (err error)
	GetAllBooks() (books []_entity.Book, err error)
	GetBookAuthors(bookId uint) (authors []_entity.Author, err error)
	GetBookById(bookId uint) (book _entity.Book, err error)
	CountBookById(bookId uint) (count uint, err error)
	UpdateBook(updatedBook _entity.Book) (book _entity.Book, err error)
	DeleteBookAuthorJunction(book _entity.Book, author _entity.Author) (err error)
	DeleteBook(bookId uint) (err error)
	AddBookToFavorite(userId uint, bookId uint) (favorite _entity.Favorite, err error)
	RemoveBookFromFavorite(userId uint, bookId uint) (err error)
	GetAllFavorites(userId uint) (favorites []_entity.Favorite, err error)
	AddBookToWishlist(userId uint, newWish _entity.Wish) (wish _entity.Wish, err error)
	RemoveBookFromWishlist(wishId uint) (err error)
	GetWishesByUserId(userId uint) (wishes []_entity.Wish, err error)
	GetWishById(userId uint, wishId uint) (wish _entity.Wish, err error)
	CreateWishAuthorJunction(wish _entity.Wish, author _entity.Author) (err error)
	GetWishAuthors(wishId uint) (authors []_entity.Author, err error)
	UpdateWish(updatedWish _entity.Wish) (wish _entity.Wish, err error)
	DeleteWishAuthorJunction(wish _entity.Wish, author _entity.Author) (err error)
}
