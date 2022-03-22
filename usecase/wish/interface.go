package wish

import (
	_model "plain-go/public-library/model"
)

type Wish interface {
	AddBookToWishlist(userId uint, req _model.AddBookToWishlistRequest) (res _model.AddBookToWishlistResponse, code int, message string)
	RemoveBookFromWishlist(userId uint, wishId uint) (code int, message string)
	GetAllWishes() (res []_model.GetWishesByUserIdResponse, code int, message string)
	GetAllWishesByUserId(userId uint) (res _model.GetWishesByUserIdResponse, code int, message string)
}
