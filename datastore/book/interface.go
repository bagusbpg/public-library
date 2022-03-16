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
	// DeleteBook(bookId int) (err error)
}
