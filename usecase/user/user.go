package user

import (
	"errors"
	"log"
	"net/http"
	_userRepository "plain-go/public-library/datastore/user"
	_entity "plain-go/public-library/entity"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repository _userRepository.User
}

func New(user _userRepository.User) *UserUseCase {
	return &UserUseCase{repository: user}
}

func (uuc UserUseCase) SignUp(req _model.SignUpRequest) (res _model.SignUpResponse, code int, err error) {
	// prepare input string
	name := strings.Title(strings.TrimSpace(req.Name))
	email := strings.TrimSpace(req.Email)
	phone := strings.TrimSpace(req.Phone)
	password := strings.TrimSpace(req.Password)

	check := []string{name, email, phone, password}

	for _, s := range check {
		// check if required input is empty
		if s == "" {
			log.Println("empty input")
			code, err = http.StatusBadRequest, errors.New("empty input")
			return
		}

		// check if there is any forbidden character in required field
		if strings.ContainsAny(strings.ReplaceAll(s, " ", ""), ";--") {
			log.Println("forbidden character")
			code, err = http.StatusBadRequest, errors.New("forbidden chacarter")
			return
		}
	}

	// check if email pattern invalid
	if err = _helper.CheckEmailPattern(email); err != nil {
		code = http.StatusBadRequest
		return
	}

	// check if phone pattern invalid
	if err = _helper.CheckPhonePattern(phone); err != nil {
		code = http.StatusBadRequest
		return
	}

	// check if password pattern invalid
	if err = _helper.CheckPasswordPattern(password); err != nil {
		code = http.StatusBadRequest
		return
	}

	// calling repository
	existingUser, code, err := uuc.repository.GetUserByEmail(email)

	// detect failure in repository
	if err != nil {
		return
	}

	// check if email is already used by other account
	if existingUser != (_entity.User{}) {
		log.Println("email already used")
		code, err = http.StatusConflict, errors.New("email already used")
		return
	}

	// hashing password before storing in database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	// detect failure in hashing password
	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("failed to hash password")
		return
	}

	// prepare input to repository
	newUser := _entity.User{}
	newUser.Name = name
	newUser.Email = email
	newUser.Phone = phone
	newUser.Password = string(hashedPassword)

	// calling repository
	res.User, code, err = uuc.repository.CreateNewUser(newUser)

	// detect failure in repository
	if err != nil {
		return
	}

	// formatting response
	res.User.Password = ""
	res.User.CreatedAt, _ = _helper.TimeFormatter(res.User.CreatedAt)
	res.User.UpdatedAt, _ = _helper.TimeFormatter(res.User.UpdatedAt)

	return
}

func (uuc UserUseCase) Login(req _model.LoginRequest) (res _model.LoginResponse, code int, err error) {
	// prepare input string
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	check := []string{email, password}

	for _, s := range check {
		// check if required input is empty
		if s == "" {
			log.Println("empty input")
			code, err = http.StatusBadRequest, errors.New("empty input")
			return
		}

		// check if there is any forbidden character in required field
		if strings.ContainsAny(strings.ReplaceAll(s, " ", ""), ";--") {
			log.Println("forbidden character")
			code, err = http.StatusBadRequest, errors.New("forbidden chacarter")
			return
		}
	}

	// calling repository
	existingUser, code, err := uuc.repository.GetUserByEmail(email)

	// detect failure in repository
	if err != nil {
		return
	}

	// check if user does not exist
	if existingUser == (_entity.User{}) {
		log.Println("user not found")
		code, err = http.StatusNotFound, errors.New("user not found")
		return
	}

	// check if password matches
	if err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(password)); err != nil {
		log.Println(err)
		code, err = http.StatusUnauthorized, errors.New("password mismatch")
		return
	}

	// create token
	token, expire, err := _helper.CreateToken(existingUser.Id, existingUser.Role)

	// detect error while creating token
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	// formatting response
	existingUser.Password = ""
	existingUser.CreatedAt = existingUser.CreatedAt.Add(7 * time.Hour)
	existingUser.UpdatedAt = existingUser.UpdatedAt.Add(7 * time.Hour)

	res.Token = token
	res.Expire = expire
	res.User = existingUser

	return
}
