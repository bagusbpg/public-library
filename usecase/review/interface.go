package review

import (
	_model "plain-go/public-library/model"
)

type Review interface {
	CreateReview(userId uint, bookId uint, req _model.CreateReviewRequest) (res _model.CreateReviewResponse, code int, message string)
}
