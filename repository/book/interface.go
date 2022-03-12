package book

import (
	_entity "plain-go/public-library/entity"
)

type Book interface {
	GetAllAuthors() (authors []_entity.Author, code int, err error)
	CreateNewBook(newBook _entity.Book) (book _entity.Book, code int, err error)
	CreateNewAuthor(newAuthor _entity.Author) (author _entity.Author, code int, err error)
	CreateBookAuthorJunction(book _entity.Book, author _entity.Author) (code int, err error)
	CreateBookItem(book _entity.Book) (code int, err error)
	GetAllBooks() (books []_entity.Book, code int, err error)
	GetBookAuthors(bookId int) (authors []_entity.Author, code int, err error)
	GetBookById(bookId int) (book _entity.Book, code int, err error)
	UpdateBook(updatedBook _entity.Book) (book _entity.Book, code int, err error)
	DeleteBook(bookId int) (code int, err error)
}
