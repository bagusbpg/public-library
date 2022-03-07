package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_common "plain-go/public-library/delivery/common"
	_entity "plain-go/public-library/entity"
	_userRepository "plain-go/public-library/repository/user"

	"golang.org/x/crypto/bcrypt"
)

type UserController struct {
	repository _userRepository.User
}

func New(user _userRepository.User) *UserController {
	return &UserController{repository: user}
}

func (uc UserController) CreateNewUser() http.HandlerFunc {
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

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
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
