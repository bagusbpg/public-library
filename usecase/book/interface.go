package book

import (
	_model "plain-go/public-library/model"
)

type Book interface {
	CreateBook(req _model.CreateBookRequest) (res _model.CreateBookResponse, code int, message string)
}
