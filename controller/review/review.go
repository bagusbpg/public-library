package review

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_mw "plain-go/public-library/app/middleware"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
	_reviewUseCase "plain-go/public-library/usecase/review"
	"strconv"
	"strings"
)

type ReviewController struct {
	usecase _reviewUseCase.Review
}

func New(review _reviewUseCase.Review) *ReviewController {
	return &ReviewController{usecase: review}
}

func (rc ReviewController) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		res, code, message := rc.usecase.GetAllReviews()

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (rc ReviewController) Create() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")
		userId, _, _ := _helper.ExtractToken(token)

		bookId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "internal server error", nil)
			return
		}

		defer r.Body.Close()

		req := _model.CreateReviewRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := rc.usecase.CreateReview(uint(userId), uint(bookId), req)

		if code != http.StatusCreated {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (rc ReviewController) GetAllByBook() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bookId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		res, code, message := rc.usecase.GetAllReviewsByBookId(uint(bookId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (rc ReviewController) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bookId, _ := strconv.Atoi(_mw.GetParam(r)[0])
		reviewId, _ := strconv.Atoi(_mw.GetParam(r)[1])

		res, code, message := rc.usecase.GetReviewByReviewId(uint(bookId), uint(reviewId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (rc ReviewController) Update() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")
		userId, role, _ := _helper.ExtractToken(token)

		bookId, _ := strconv.Atoi(_mw.GetParam(r)[0])
		reviewId, _ := strconv.Atoi(_mw.GetParam(r)[1])

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "internal server error", nil)
			return
		}

		defer r.Body.Close()

		req := _model.UpdateReviewRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		if role == "Member" {
			res, code, message := rc.usecase.UpdateReview(uint(userId), uint(bookId), uint(reviewId), req)

			if code != http.StatusOK {
				_model.CreateResponse(rw, code, message, nil)
				return
			}

			_model.CreateResponse(rw, code, message, res)
		} else if role == "Librarian" {
			code, message := rc.usecase.UpdateStatus(uint(bookId), uint(reviewId), req)

			_model.CreateResponse(rw, code, message, nil)
		}
	}
}

func (rc ReviewController) Delete() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")
		userId, _, _ := _helper.ExtractToken(token)

		bookId, _ := strconv.Atoi(_mw.GetParam(r)[0])
		reviewId, _ := strconv.Atoi(_mw.GetParam(r)[1])

		code, message := rc.usecase.DeleteReview(uint(userId), uint(bookId), uint(reviewId))

		_model.CreateResponse(rw, code, message, nil)
	}
}
