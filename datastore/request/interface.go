package request

import (
	_entity "plain-go/public-library/entity"
)

type Request interface {
	CountActiveRequestByUserId(userId uint) (count uint, err error)
	GetRequestByUserIdAndBookId(userId uint, bookId uint) (requests []_entity.Request, err error)
	CreateNewRequest(newRequest _entity.Request) (request _entity.Request, err error)
	GetRequestByRequestId(requestId uint) (request _entity.Request, err error)
	GetAllRequests() (requests []_entity.Request, err error)
	GetRequestStatusId(statusDesc string) (statusId uint, err error)
	Update(updatedRequest _entity.Request) (request _entity.Request, err error)
}
