package request

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_mw "plain-go/public-library/app/middleware"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
	_requestUseCase "plain-go/public-library/usecase/request"
	"strconv"
	"strings"
)

type RequestController struct {
	usecase _requestUseCase.Request
}

func New(request _requestUseCase.Request) *RequestController {
	return &RequestController{usecase: request}
}

func (rc RequestController) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		res, code, message := rc.usecase.GetAllRequests()

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (rc RequestController) GetAllByUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		res, code, message := rc.usecase.GetAllRequestsByUserId(uint(userId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (rc RequestController) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])
		requestId, _ := strconv.Atoi(_mw.GetParam(r)[1])

		res, code, message := rc.usecase.GetRequestById(uint(userId), uint(requestId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (rc RequestController) Create() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

		req := _model.CreateRequestRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := rc.usecase.CreateRequest(uint(userId), req)

		if code != http.StatusCreated {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (rc RequestController) Update() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		token := strings.TrimPrefix(r.Header.Get("authorization"), "Bearer ")
		userId, role, _ := _helper.ExtractToken(token)
		requestId, _ := strconv.Atoi(_mw.GetParam(r)[1])

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

		req := _model.UpdateRequestRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := rc.usecase.UpdateRequest(uint(userId), uint(requestId), role, req)

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}
