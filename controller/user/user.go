package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_model "plain-go/public-library/model"
	_userUseCase "plain-go/public-library/usecase/user"
	"strconv"
	"strings"
)

type UserController struct {
	usecase _userUseCase.User
}

func New(user _userUseCase.User) *UserController {
	return &UserController{usecase: user}
}

func (uc UserController) SignUp() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

		if contentType := r.Header.Get("content-type"); contentType != "application/json" {
			log.Println("unsupported content type")
			_model.CreateResponse(rw, http.StatusUnsupportedMediaType, "unsupported content type", nil)
			return
		}

		req := _model.SignUpRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := uc.usecase.SignUp(req)

		if code != http.StatusCreated {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (uc UserController) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

		if contentType := r.Header.Get("content-type"); contentType != "application/json" {
			log.Println("unsupported content type")
			_model.CreateResponse(rw, http.StatusUnsupportedMediaType, "unsupported content type", nil)
			return
		}

		req := _model.LoginRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := uc.usecase.Login(req)

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (uc UserController) GetUpdateDelete() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(strings.SplitAfter(r.URL.Path, "/")[2])

		existing, code, message := uc.usecase.GetUserById(uint(userId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		switch r.Method {
		case http.MethodGet:
			existing.User.Password = ""
			_model.CreateResponse(rw, code, message, existing.User)
		case http.MethodPut:
			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				log.Println(err)
				_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
				return
			}

			defer r.Body.Close()

			if contentType := r.Header.Get("content-type"); contentType != "application/json" {
				log.Println("unsupported content type")
				_model.CreateResponse(rw, http.StatusUnsupportedMediaType, "unsupported content type", nil)
				return
			}

			req := _model.UpdateUserRequest{}

			if err = json.Unmarshal(body, &req); err != nil {
				log.Println(err)
				_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
				return
			}

			res, code, message := uc.usecase.UpdateUser(req, existing.User)

			if code != http.StatusOK {
				_model.CreateResponse(rw, code, message, nil)
				return
			}

			_model.CreateResponse(rw, code, message, res)
		case http.MethodDelete:
			code, message := uc.usecase.DeleteUser(uint(userId))

			_model.CreateResponse(rw, code, message, nil)
		default:
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
		}
	}
}
