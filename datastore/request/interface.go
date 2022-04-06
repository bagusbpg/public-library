package request

import (
	_entity "plain-go/public-library/entity"
)

type Request interface {
	GetAllRequests() (requests []_entity.Request, err error)
	GetAllRequestsByUserId(userId uint) (requests []_entity.SimplifiedRequest, err error)
	CountActiveRequestByUserId(userId uint) (count uint, err error)
	GetRequestByUserIdAndBookId(userId uint, bookId uint) (requests []_entity.Request, err error)
	CreateNewRequest(newRequest _entity.Request) (request _entity.Request, err error)
	GetRequestById(requestId uint) (request _entity.Request, err error)
	GetRequestStatusId(statusDesc string) (statusId uint, err error)
	Update(updatedRequest _entity.Request) (request _entity.Request, err error)
}
