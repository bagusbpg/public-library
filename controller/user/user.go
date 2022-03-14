package user

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_model "plain-go/public-library/model"
	_userUseCase "plain-go/public-library/usecase/user"
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

		contentType := r.Header.Get("content-type")

		if contentType != "application/json" {
			log.Println("unsupported content type")
			_model.CreateResponse(rw, http.StatusUnsupportedMediaType, "unsupported content type", nil)
			return
		}

		req := _model.SignUpRequest{}
		err = json.Unmarshal(body, &req)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, err := uc.usecase.SignUp(req)

		if err != nil {
			_model.CreateResponse(rw, code, err.Error(), nil)
			return
		}

		_model.CreateResponse(rw, http.StatusCreated, "success create user", res)
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

		contentType := r.Header.Get("content-type")

		if contentType != "application/json" {
			log.Println("unsupported content type")
			_model.CreateResponse(rw, http.StatusUnsupportedMediaType, "unsupported content type", nil)
			return
		}

		req := _model.LoginRequest{}
		err = json.Unmarshal(body, &req)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, err := uc.usecase.Login(req)

		if err != nil {
			_model.CreateResponse(rw, code, err.Error(), nil)
			return
		}

		_model.CreateResponse(rw, http.StatusOK, "success login", res)
	}
}

// func (uc UserController) GetUpdateDelete() http.HandlerFunc {
// 	return func(rw http.ResponseWriter, r *http.Request) {
// 		userId, _ := strconv.Atoi(strings.SplitAfter(r.URL.Path, "/")[2])

// 		existingUser, code, err := uc.repository.GetUserById(userId)

// 		if err != nil {
// 			_model.CreateResponse(rw, code, err.Error(), nil)
// 			return
// 		}

// 		if existingUser == (_entity.User{}) {
// 			_model.CreateResponse(rw, http.StatusNotFound, "user not found", nil)
// 			return
// 		}

// 		switch r.Method {
// 		case http.MethodGet:
// 			existingUser.Id = userId
// 			existingUser.Password = ""
// 			existingUser.CreatedAt = existingUser.CreatedAt.Add(7 * time.Hour)
// 			existingUser.UpdatedAt = existingUser.UpdatedAt.Add(7 * time.Hour)

// 			_model.CreateResponse(rw, http.StatusOK, "success get user", existingUser)
// 		case http.MethodPut:
// 			body, err := ioutil.ReadAll(r.Body)

// 			if err != nil {
// 				log.Println(err)
// 				_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
// 				return
// 			}

// 			defer r.Body.Close()

// 			contentType := r.Header.Get("content-type")

// 			if contentType != "application/json" {
// 				log.Println("unsupported content type")
// 				_model.CreateResponse(rw, http.StatusUnsupportedMediaType, "unsupported content type", nil)
// 				return
// 			}

// 			updatedUser := _entity.User{}
// 			err = json.Unmarshal(body, &updatedUser)

// 			if err != nil {
// 				log.Println(err)
// 				_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
// 				return
// 			}

// 			name := strings.Title(strings.TrimSpace(updatedUser.Name))
// 			email := strings.TrimSpace(updatedUser.Email)
// 			phone := strings.TrimSpace(updatedUser.Phone)
// 			password := strings.TrimSpace(updatedUser.Password)

// 			check := []string{name, email, phone, password}
// 			flag := true

// 			for _, s := range check {
// 				if strings.ContainsAny(strings.ReplaceAll(s, " ", ""), ";--") {
// 					log.Println("forbidden character")
// 					_model.CreateResponse(rw, http.StatusBadRequest, "forbidden chacarter", nil)
// 					return
// 				}

// 				if s != "" {
// 					flag = false
// 				}
// 			}

// 			if flag {
// 				log.Println("no update was performed")
// 				_model.CreateResponse(rw, http.StatusBadRequest, "no update was performed", nil)
// 				return
// 			}

// 			if name != "" {
// 				existingUser.Name = name
// 			}

// 			if email != "" {
// 				if err = _helper.CheckEmailPattern(email); err != nil {
// 					_model.CreateResponse(rw, http.StatusBadRequest, err.Error(), nil)
// 					return
// 				}

// 				existingUser.Email = email
// 			}

// 			if phone != "" {
// 				if err = _helper.CheckPhonePattern(phone); err != nil {
// 					_model.CreateResponse(rw, http.StatusBadRequest, err.Error(), nil)
// 					return
// 				}

// 				existingUser.Phone = phone
// 			}

// 			if password != "" {
// 				if err = _helper.CheckPasswordPattern(password); err != nil {
// 					_model.CreateResponse(rw, http.StatusBadRequest, err.Error(), nil)
// 					return
// 				}

// 				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

// 				if err != nil {
// 					log.Println(err)
// 					_model.CreateResponse(rw, http.StatusInternalServerError, "failed to hash password", nil)
// 					return
// 				}

// 				existingUser.Password = string(hashedPassword)
// 			}

// 			user, code, err := uc.repository.UpdateUser(existingUser)
// 			user.Id = userId
// 			user.Password = ""
// 			user.CreatedAt = user.CreatedAt.Add(7 * time.Hour)
// 			user.UpdatedAt, _ = _helper.TimeFormatter(user.UpdatedAt)

// 			if err != nil {
// 				_model.CreateResponse(rw, code, err.Error(), nil)
// 				return
// 			}

// 			_model.CreateResponse(rw, http.StatusCreated, "success update user", user)
// 		case http.MethodDelete:
// 			code, err := uc.repository.DeleteUser(userId)

// 			if err != nil {
// 				_model.CreateResponse(rw, code, err.Error(), nil)
// 				return
// 			}

// 			_model.CreateResponse(rw, http.StatusOK, "success delete user", nil)
// 		default:
// 			log.Println("method not allowed")
// 			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
// 			return
// 		}
// 	}
// }
