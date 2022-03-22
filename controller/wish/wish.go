package wish

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_mw "plain-go/public-library/app/middleware"
	_model "plain-go/public-library/model"
	_wishUseCase "plain-go/public-library/usecase/wish"
	"strconv"
)

type WishController struct {
	usecase _wishUseCase.Wish
}

func New(wish _wishUseCase.Wish) *WishController {
	return &WishController{usecase: wish}
}

func (wc WishController) AddBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

		req := _model.AddBookToWishlistRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := wc.usecase.AddBookToWishlist(uint(userId), req)

		if code != http.StatusCreated {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (wc WishController) RemoveBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])
		wishId, _ := strconv.Atoi(_mw.GetParam(r)[1])

		code, message := wc.usecase.RemoveBookFromWishlist(uint(userId), uint(wishId))

		_model.CreateResponse(rw, code, message, nil)
	}
}

func (wc WishController) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		res, code, message := wc.usecase.GetAllWishes()

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (wc WishController) GetAllByUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		res, code, message := wc.usecase.GetAllWishesByUserId(uint(userId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
		}

		_model.CreateResponse(rw, code, message, res.Wishes)
	}
}

func (wc WishController) Update() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])
		wishId, _ := strconv.Atoi(_mw.GetParam(r)[1])

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

		req := _model.UpdateWishRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := wc.usecase.UpdateWish(req, uint(userId), uint(wishId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}
