package user

import (
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

func (uuc UserUseCase) SignUp(req _model.SignUpRequest) (res _model.SignUpResponse, code int, message string) {
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
			code, message = http.StatusBadRequest, "empty input"
			return
		}

		// check if there is any forbidden character in required field
		if strings.Contains(strings.ReplaceAll(s, " ", ""), ";--") {
			log.Println("forbidden character")
			code, message = http.StatusBadRequest, "forbidden chacarter"
			return
		}
	}

	// check if email pattern invalid
	if err := _helper.CheckEmailPattern(email); err != nil {
		code, message = http.StatusBadRequest, err.Error()
		return
	}

	// check if phone pattern invalid
	if err := _helper.CheckPhonePattern(phone); err != nil {
		code, message = http.StatusBadRequest, err.Error()
		return
	}

	// check if password pattern invalid
	if err := _helper.CheckPasswordPattern(password); err != nil {
		code, message = http.StatusBadRequest, err.Error()
		return
	}

	// check if email is already used by other account
	newUser, err := uuc.repository.GetUserByEmail(email)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if newUser.Email != "" {
		log.Println("email already used")
		code, message = http.StatusConflict, "email already used"
		return
	}

	// hashing password before storing in database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	// detect failure in hashing password
	if err != nil {
		log.Println(err)
		code, message = http.StatusInternalServerError, "failed to hash password"
		return
	}

	// prepare input to repository
	now := time.Now()
	newUser.Role = "Member"
	newUser.Name = name
	newUser.Email = email
	newUser.Phone = phone
	newUser.Password = string(hashedPassword)
	newUser.CreatedAt = now
	newUser.UpdatedAt = now

	// calling repository
	res.User, err = uuc.repository.CreateNewUser(newUser)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// formatting response
	res.User.Password = ""
	res.User.CreatedAt, _ = _helper.TimeFormatter(res.User.CreatedAt)
	res.User.UpdatedAt, _ = _helper.TimeFormatter(res.User.UpdatedAt)
	code, message = http.StatusCreated, "success create user"

	return
}

func (uuc UserUseCase) Login(req _model.LoginRequest) (res _model.LoginResponse, code int, message string) {
	// prepare input string
	email := strings.TrimSpace(req.Email)
	password := strings.TrimSpace(req.Password)

	check := []string{email, password}

	for _, s := range check {
		// check if required input is empty
		if s == "" {
			log.Println("empty input")
			code, message = http.StatusBadRequest, "empty input"
			return
		}

		// check if there is any forbidden character in required field
		if strings.Contains(strings.ReplaceAll(s, " ", ""), ";--") {
			log.Println("forbidden character")
			code, message = http.StatusBadRequest, "forbidden chacarter"
			return
		}
	}

	// calling repository
	user, err := uuc.repository.GetUserByEmail(email)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	res.User = user

	// check if user does not exist
	if res.User == (_entity.User{}) {
		log.Println("user not found")
		code, message = http.StatusNotFound, "user not found"
		return
	}

	// check if password matches
	if err := bcrypt.CompareHashAndPassword([]byte(res.User.Password), []byte(password)); err != nil {
		log.Println(res.User.Password, password, err)
		code, message = http.StatusUnauthorized, "password mismatch"
		return
	}

	// create token
	token, expire, err := _helper.CreateToken(res.User.Id, res.User.Role)

	// detect error while creating token
	if err != nil {
		code, message = http.StatusInternalServerError, "failed to create token"
		return
	}

	// formatting response
	res.User.Password = ""
	res.User.CreatedAt = res.User.CreatedAt.Add(7 * time.Hour)
	res.User.UpdatedAt = res.User.UpdatedAt.Add(7 * time.Hour)
	res.Token = token
	res.Expire = expire
	code, message = http.StatusOK, "success login"

	return
}

func (uuc UserUseCase) GetAllUsers() (res _model.GetAllUsersResponse, code int, message string) {
	// calling repository
	users, err := uuc.repository.GetAllUsers()

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, user := range users {
		// formatting response
		user.CreatedAt, _ = _helper.TimeFormatter(user.CreatedAt)
		user.UpdatedAt, _ = _helper.TimeFormatter(user.UpdatedAt)
		res.Users = append(res.Users, user)
	}

	code, message = http.StatusOK, "success get all users"

	return
}

func (uuc UserUseCase) GetUserById(userId uint) (res _model.GetUserByIdResponse, code int, message string) {
	// calling repository
	user, err := uuc.repository.GetUserById(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	res.User = user

	// check if user does not exist
	if res.User.Name == "" {
		log.Println("user not found")
		code, message = http.StatusNotFound, "user not found"
		return
	}

	// formatting response
	res.User.CreatedAt = res.User.CreatedAt.Add(7 * time.Hour)
	res.User.UpdatedAt = res.User.UpdatedAt.Add(7 * time.Hour)
	code, message = http.StatusOK, "success get user"

	return
}

func (uuc UserUseCase) UpdateUser(req _model.UpdateUserRequest, userId uint) (res _model.UpdateUserResponse, code int, message string) {
	// check user existence
	user, err := uuc.repository.GetUserById(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// check if user does not exist
	if user.Name == "" {
		log.Println("user not found")
		code, message = http.StatusNotFound, "user not found"
		return
	}

	// prepare input string
	name := strings.Title(strings.TrimSpace(req.Name))
	email := strings.TrimSpace(req.Email)
	phone := strings.TrimSpace(req.Phone)
	password := strings.TrimSpace(req.Password)

	check := []string{name, email, phone, password}
	flag := true

	for _, s := range check {
		// check if there is any forbidden character
		if strings.Contains(strings.ReplaceAll(s, " ", ""), ";--") {
			log.Println("forbidden character")
			code, message = http.StatusBadRequest, "forbidden chacarter"
			return
		}
	}

	if name != "" && name != user.Name {
		user.Name = name
		flag = false
	}

	if email != "" && email != user.Email {
		// check if email pattern invalid
		if err := _helper.CheckEmailPattern(email); err != nil {
			code, message = http.StatusBadRequest, err.Error()
			return
		}

		user.Email = email
		flag = false
	}

	if phone != "" && phone != user.Phone {
		// check if phone pattern invalid
		if err := _helper.CheckPhonePattern(phone); err != nil {
			code, message = http.StatusBadRequest, err.Error()
			return
		}

		user.Phone = phone
		flag = false
	}

	if password != "" {
		// check if password pattern invalid
		if err := _helper.CheckPasswordPattern(password); err != nil {
			code, message = http.StatusBadRequest, err.Error()
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			// hashing password before storing in database
			hashedPassword, errhash := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

			// detect failure in hashing password
			if errhash != nil {
				log.Println(errhash)
				code, message = http.StatusInternalServerError, "failed to hash password"
				return
			}

			user.Password = string(hashedPassword)
			flag = false
		}
	}

	// check if no field is updated
	if flag {
		log.Println("no update was performed")
		code, message = http.StatusBadRequest, "no update was performed"
		return
	}

	// calling respository
	user.UpdatedAt = time.Now()
	_user, err := uuc.repository.UpdateUser(user)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// formatting response
	res.User = _user
	res.User.Id = user.Id
	res.User.Password = ""
	res.User.UpdatedAt, _ = _helper.TimeFormatter(res.User.UpdatedAt)
	code, message = http.StatusOK, "success update user"

	return
}

func (uuc UserUseCase) DeleteUser(userId uint) (code int, message string) {
	// calling repository
	err := uuc.repository.DeleteUser(userId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	code, message = http.StatusOK, "success delete user"

	return
}
