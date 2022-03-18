package favorite

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_helper "plain-go/public-library/helper"
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

func (fc FavoriteController) AddBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")

		loginId, _, _ := _helper.ExtractToken(token)

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

		code, message := fc.usecase.AddBookToFavorite(uint(loginId), req.BookId)

		_model.CreateResponse(rw, code, message, nil)
	}
}

func (fc FavoriteController) RemoveBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		str := strings.SplitAfter(r.URL.Path, "/")
		userId, _ := strconv.Atoi(str[len(str)-1])

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

func (fc FavoriteController) GetAll() http.HandlerFunc {
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
