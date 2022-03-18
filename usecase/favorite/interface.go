package favorite

import (
	_model "plain-go/public-library/model"
)

type Favorite interface {
	AddBookToFavorite(userId uint, bookId uint) (code int, message string)
	GetAllFavorites(userId uint) (res _model.GetAllFavoritesResponse, code int, message string)
}
