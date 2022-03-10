package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_common "plain-go/public-library/delivery/common"
	_helper "plain-go/public-library/delivery/helper"
	_entity "plain-go/public-library/entity"
	_userRepository "plain-go/public-library/repository/user"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repository _userRepository.User
}

func New(user _userRepository.User) *UserController {
	return &UserController{repository: user}
}

func (uc UserController) SignUp() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			log.Println("method not allowed")
			_common.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_common.CreateResponse(rw, http.StatusInternalServerError, "internal server error", nil)
			return
		}

		defer r.Body.Close()

		contentType := r.Header.Get("content-type")

		if contentType != "application/json" {
			log.Println("unsupported content type")
			_common.CreateResponse(rw, http.StatusUnsupportedMediaType, "unsupported content type", nil)
			return
		}

		newUser := _entity.User{}
		err = json.Unmarshal(body, &newUser)

		if err != nil {
			log.Println(err)
			_common.CreateResponse(rw, http.StatusBadRequest, "failed in binding", nil)
			return
		}

		name := strings.Title(strings.TrimSpace(newUser.Name))
		email := strings.TrimSpace(newUser.Email)
		phone := strings.TrimSpace(newUser.Phone)
		password := strings.TrimSpace(newUser.Password)

		check := []string{name, email, phone, password}

		for _, s := range check {
			if s == "" {
				log.Println("empty input")
				_common.CreateResponse(rw, http.StatusBadRequest, "empty input", nil)
				return
			}

			if strings.ContainsAny(strings.ReplaceAll(s, " ", ""), ";--") {
				log.Println("forbidden character")
				_common.CreateResponse(rw, http.StatusBadRequest, "forbidden chacarter", nil)
				return
			}
		}

		if err = _helper.CheckEmailPattern(email); err != nil {
			_common.CreateResponse(rw, http.StatusBadRequest, err.Error(), nil)
			return
		}

		if err = _helper.CheckPhonePattern(phone); err != nil {
			_common.CreateResponse(rw, http.StatusBadRequest, err.Error(), nil)
			return
		}

		if err = _helper.CheckPasswordPattern(password); err != nil {
			_common.CreateResponse(rw, http.StatusBadRequest, err.Error(), nil)
			return
		}

		existingUser, code, err := uc.repository.GetUserByEmail(newUser.Email)

		if err != nil {
			_common.CreateResponse(rw, code, err.Error(), nil)
			return
		}

		if existingUser != (_entity.User{}) {
			_common.CreateResponse(rw, http.StatusConflict, "email already used", nil)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

		if err != nil {
			log.Println(err)
			_common.CreateResponse(rw, http.StatusInternalServerError, "internal server error", nil)
			return
		}

		newUser.Password = string(hashedPassword)

		user, code, err := uc.repository.CreateNewUser(newUser)
		user.Password = ""

		if err != nil {
			_common.CreateResponse(rw, code, err.Error(), nil)
			return
		}

		_common.CreateResponse(rw, http.StatusCreated, "success create user", user)
	}
}

func (uc UserController) Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			log.Println("method not allowed")
			_common.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_common.CreateResponse(rw, http.StatusInternalServerError, "internal server error", nil)
			return
		}

		defer r.Body.Close()

		contentType := r.Header.Get("content-type")

		if contentType != "application/json" {
			log.Println("unsupported content type")
			_common.CreateResponse(rw, http.StatusUnsupportedMediaType, "unsupported content type", nil)
			return
		}

		loginUser := _entity.User{}
		err = json.Unmarshal(body, &loginUser)

		if err != nil {
			log.Println(err)
			_common.CreateResponse(rw, http.StatusBadRequest, "failed in binding", nil)
			return
		}

		email := strings.TrimSpace(loginUser.Email)
		password := strings.TrimSpace(loginUser.Password)

		check := []string{email, password}

		for _, s := range check {
			if s == "" {
				log.Println("empty input")
				_common.CreateResponse(rw, http.StatusBadRequest, "empty input", nil)
				return
			}

			if strings.ContainsAny(strings.ReplaceAll(s, " ", ""), ";--") {
				log.Println("forbidden character")
				_common.CreateResponse(rw, http.StatusBadRequest, "forbidden chacarter", nil)
				return
			}
		}

		existingUser, code, err := uc.repository.GetUserByEmail(loginUser.Email)

		if err != nil {
			_common.CreateResponse(rw, code, err.Error(), nil)
			return
		}

		if existingUser == (_entity.User{}) {
			_common.CreateResponse(rw, http.StatusNotFound, "user not found", nil)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
			fmt.Println(existingUser.Password, password)
			log.Println(err)
			_common.CreateResponse(rw, http.StatusUnauthorized, "password mismatch", nil)
			return
		}

		token, expire, err := _helper.CreateToken(existingUser.Id, existingUser.Role)

		if err != nil {
			_common.CreateResponse(rw, http.StatusInternalServerError, "internal server error", nil)
			return
		}

		existingUser.Password = ""
		existingUser.CreatedAt = existingUser.CreatedAt.Add(7 * time.Hour)
		existingUser.UpdatedAt = existingUser.UpdatedAt.Add(7 * time.Hour)

		_common.CreateResponse(rw, http.StatusOK, "success login", map[string]interface{}{"user": existingUser, "token": token, "expire": expire})
	}
}
