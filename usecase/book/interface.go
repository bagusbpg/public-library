package book

import (
	_model "plain-go/public-library/model"
)

type Book interface {
	CreateBook(req _model.CreateBookRequest) (res _model.CreateBookResponse, code int, message string)
	GetAllBooks(params _model.GetAllBooksRequest) (res _model.GetAllBooksResponse, code int, message string)
	GetBookById(bookId uint) (res _model.GetBookByIdResponse, code int, message string)
	UpdateBook(req _model.UpdateBookRequest, bookId uint) (res _model.UpdateBookResponse, code int, message string)
	DeleteBook(bookId uint) (code int, message string)
}
