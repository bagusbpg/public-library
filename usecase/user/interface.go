package user

import (
	_model "plain-go/public-library/model"
)

type User interface {
	SignUp(req _model.SignUpRequest) (res _model.SignUpResponse, code int, message string)
	Login(req _model.LoginRequest) (res _model.LoginResponse, code int, message string)
	GetAllUsers() (res _model.GetAllUsersResponse, code int, message string)
	GetUserById(userId uint) (res _model.GetUserByIdResponse, code int, message string)
	UpdateUser(req _model.UpdateUserRequest, userId uint) (res _model.UpdateUserResponse, code int, message string)
	DeleteUser(userId uint) (code int, message string)
}
