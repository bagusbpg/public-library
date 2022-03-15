package model

import (
	_entity "plain-go/public-library/entity"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token  string       `json:"token"`
	Expire int64        `json:"expire"`
	User   _entity.User `json:"user"`
}

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	User _entity.User `json:"user"`
}

type GetUserByIdResponse struct {
	User _entity.User `json:"user"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type UpdateUserResponse struct {
	User _entity.User `json:"user"`
}
