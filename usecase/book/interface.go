package book

import (
	_entity "plain-go/public-library/entity"
	_model "plain-go/public-library/model"
)

type Book interface {
	CreateBook(req _model.CreateBookRequest) (res _model.CreateBookResponse, code int, message string)
	GetAllBooks() (res _model.GetAllBooksResponse, code int, message string)
	GetBookById(bookId uint) (res _model.GetBookByIdResponse, code int, message string)
	UpdateBook(req _model.UpdateBookRequest, book _entity.Book) (res _model.UpdateBookResponse, code int, message string)
}
