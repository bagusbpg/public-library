package review

import (
	_model "plain-go/public-library/model"
)

type Review interface {
	CreateReview(userId uint, bookId uint, req _model.CreateReviewRequest) (res _model.CreateReviewResponse, code int, message string)
	GetReviewByReviewId(reviewId uint) (res _model.GetReviewByReviewIdResponse, code int, message string)
	GetAllReviewsByBookId(bookId uint) (res _model.GetAllReviewsByBookIdResponse, code int, message string)
	UpdateReview(userId uint, reviewId uint, req _model.UpdateReviewRequest) (res _model.UpdateReviewResponse, code int, message string)
	DeleteReview(userId uint, reviewId uint) (code int, message string)
}
