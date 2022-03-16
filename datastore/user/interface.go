package user

import (
	_entity "plain-go/public-library/entity"
)

type User interface {
	GetUserByEmail(email string) (user _entity.User, err error)
	CreateNewUser(newUser _entity.User) (user _entity.User, err error)
	GetAllUsers() (users []_entity.User, err error)
	GetUserById(userId uint) (user _entity.User, err error)
	UpdateUser(updatedUser _entity.User) (user _entity.User, err error)
	DeleteUser(userId uint) (err error)
}
