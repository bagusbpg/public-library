package user

import (
	_entity "plain-go/public-library/entity"
)

type User interface {
	CreateNewUser(newUser _entity.User) (user _entity.User, code int, err error)
	GetUserByEmail(email string) (user _entity.User, code int, err error)
	GetUserById(userId int) (user _entity.User, code int, err error)
	UpdateUser(updatedUser _entity.User) (user _entity.User, code int, err error)
	DeleteUser(userId int) (code int, err error)
}
