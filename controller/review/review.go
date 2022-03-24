package review

import (
	_reviewUseCase "plain-go/public-library/usecase/review"
)

type ReviewController struct {
	usecase _reviewUseCase.Review
}

func New(review _reviewUseCase.Review) *ReviewController {
	return &ReviewController{usecase: review}
}
