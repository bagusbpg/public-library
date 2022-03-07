package user

import (
	_entity "plain-go/public-library/entity"
)

type User interface {
	CreateNewUser(newUser _entity.User) (user _entity.User, code int, err error)
	GetUserDetail(userId int) (user _entity.User, code int, err error)
	UpdateUserDetail(updatedUser _entity.User) (user _entity.User, code int, err error)
	DeleteUser(userId int) (code int, err error)
}
