package user

import (
	_entity "plain-go/public-library/entity"
	_model "plain-go/public-library/model"
)

type User interface {
	SignUp(req _model.SignUpRequest) (res _model.SignUpResponse, code int, message string)
	Login(req _model.LoginRequest) (res _model.LoginResponse, code int, message string)
	GetUserById(userId int) (res _model.GetUserByIdResponse, code int, message string)
	UpdateUser(req _model.UpdateUserRequest, user _entity.User) (res _model.UpdateUserResponse, code int, message string)
	DeleteUser(userId int) (code int, message string)
}
