package favorite

import (
	_model "plain-go/public-library/model"
)

type Favorite interface {
	GetAllFavorites(userId uint) (res _model.GetAllFavoritesResponse, code int, message string)
}
