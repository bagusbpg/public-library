package favorite

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_mw "plain-go/public-library/app/middleware"
	_model "plain-go/public-library/model"
	_favoriteUseCase "plain-go/public-library/usecase/favorite"
	"strconv"
)

type FavoriteController struct {
	usecase _favoriteUseCase.Favorite
}

func New(favorite _favoriteUseCase.Favorite) *FavoriteController {
	return &FavoriteController{usecase: favorite}
}

func (fc FavoriteController) AddBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

		req := _model.AddBookToFavoriteRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := fc.usecase.AddBookToFavorite(uint(userId), req.BookId)

		if code != http.StatusCreated {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (fc FavoriteController) RemoveBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

		req := _model.RemoveBookFromFavoriteRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		code, message := fc.usecase.RemoveBookFromFavorite(uint(userId), req.BookId)

		_model.CreateResponse(rw, code, message, nil)
	}
}

func (fc FavoriteController) GetAllByUserId() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		res, code, message := fc.usecase.GetAllFavoritesByUserId(uint(userId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}
