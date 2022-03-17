package favorite

import (
	"net/http"
	_model "plain-go/public-library/model"
	_favoriteUseCase "plain-go/public-library/usecase/favorite"
	"strconv"
	"strings"
)

type FavoriteController struct {
	usecase _favoriteUseCase.Favorite
}

func New(favorite _favoriteUseCase.Favorite) *FavoriteController {
	return &FavoriteController{usecase: favorite}
}

func (fc FavoriteController) GetAllFavorites() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		str := strings.SplitAfter(r.URL.Path, "/")
		userId, _ := strconv.Atoi(str[len(str)-1])

		res, code, message := fc.usecase.GetAllFavorites(uint(userId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}
