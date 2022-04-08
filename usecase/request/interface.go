package request

import (
	_model "plain-go/public-library/model"
)

type Request interface {
	GetAllRequests() (res _model.GetAllRequestResponse, code int, message string)
	GetAllRequestsByUserId(userId uint) (res _model.GetAllRequestByUserIdResponse, code int, message string)
	GetRequestById(userId uint, requestId uint) (res _model.GetRequestByIdResponse, code int, message string)
	CreateRequest(userId uint, req _model.CreateRequestRequest) (res _model.CreateRequestResponse, code int, message string)
	UpdateRequest(userId uint, requestId uint, role string, req _model.UpdateRequestRequest) (res _model.UpdateRequestResponse, code int, message string)
}
