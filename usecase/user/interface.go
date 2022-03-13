package user

import _model "plain-go/public-library/model"

type User interface {
	SignUp(req _model.SignUpRequest) (res _model.SignUpResponse, code int, err error)
	Login(req _model.LoginRequest) (res _model.LoginResponse, code int, err error)
}
