package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_mw "plain-go/public-library/app/middleware"
	_model "plain-go/public-library/model"
	_userUseCase "plain-go/public-library/usecase/user"
	"strconv"
)

type UserController struct {
	usecase _userUseCase.User
}

func New(user _userUseCase.User) *UserController {
	return &UserController{usecase: user}
}

func (uc UserController) SignUp() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

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

func (uc UserController) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		res, code, message := uc.usecase.GetAllUsers()

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (uc UserController) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

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

func (uc UserController) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		res, code, message := uc.usecase.GetUserById(uint(userId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		res.User.Password = ""

		_model.CreateResponse(rw, code, message, res.User)
	}
}

func (uc UserController) Update() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

		req := _model.UpdateUserRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := uc.usecase.UpdateUser(req, uint(userId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (uc UserController) Delete() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		_, code, message := uc.usecase.GetUserById(uint(userId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		code, message = uc.usecase.DeleteUser(uint(userId))

		_model.CreateResponse(rw, code, message, nil)
	}
}
