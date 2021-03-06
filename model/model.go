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

type GetAllUsersResponse struct {
	Users []_entity.User `json:"users"`
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

type CreateBookRequest struct {
	Title       string                `json:"title"`
	Author      []CreateAuthorRequest `json:"author"`
	Publisher   string                `json:"publisher"`
	Language    string                `json:"language"`
	Pages       uint                  `json:"pages"`
	Category    string                `json:"category"`
	ISBN13      string                `json:"isbn13"`
	Description string                `json:"description"`
	Quantity    uint                  `json:"quantity"`
}

type CreateBookResponse struct {
	Book _entity.Book `json:"book"`
}

type CreateAuthorRequest struct {
	Name string `json:"name"`
}

type GetAllBooksRequest struct {
	Page     int
	Records  int
	Category string
	Keyword  string
	SortBy   string
	SortMode string
}

type GetAllBooksResponse struct {
	Books []_entity.Book `json:"books"`
	Count uint           `json:"count"`
}

type GetBookByIdResponse struct {
	Book _entity.Book `json:"book"`
}

type UpdateBookRequest struct {
	Title       string                `json:"title"`
	Author      []CreateAuthorRequest `json:"author"`
	Publisher   string                `json:"publisher"`
	Language    string                `json:"language"`
	Pages       uint                  `json:"pages"`
	Category    string                `json:"category"`
	ISBN13      string                `json:"isbn13"`
	Description string                `json:"description"`
}

type UpdateBookResponse struct {
	Book _entity.Book `json:"book"`
}

type AddBookToFavoriteRequest struct {
	BookId uint `json:"book_id"`
}

type AddBookToFavoriteResponse struct {
	Favorite _entity.Favorite `json:"favorite"`
}

type RemoveBookFromFavoriteRequest struct {
	BookId uint `json:"book_id"`
}

type GetAllFavoritesResponse struct {
	User      _entity.User       `json:"user"`
	Favorites []_entity.Favorite `json:"favorites"`
}

type GetAllWishes struct {
	Wishes []_entity.SimplifiedWish `json:"wishes"`
}

type GetWishesByUserIdResponse struct {
	User   _entity.User   `json:"user"`
	Wishes []_entity.Wish `json:"wishes"`
}

type AddBookToWishlistRequest struct {
	Title    string                `json:"title"`
	Author   []CreateAuthorRequest `json:"author"`
	Category string                `json:"category"`
	Note     string                `json:"note"`
}

type AddBookToWishlistResponse struct {
	Wish _entity.Wish `json:"wish"`
}

type UpdateWishRequest struct {
	Title    string                `json:"title"`
	Author   []CreateAuthorRequest `json:"authors"`
	Category string                `json:"category"`
	Note     string                `json:"note"`
}

type UpdateWishResponse struct {
	Wish _entity.Wish `json:"wish"`
}

type CreateReviewRequest struct {
	Star    uint   `json:"star"`
	Content string `json:"content"`
}

type GetAllReviewsResponse struct {
	Reviews []_entity.SimplifiedReview `json:"reviews"`
}

type GetReviewByIdResponse struct {
	Review _entity.SimplifiedReview `json:"review"`
}

type GetAllReviewsByBookIdResponse struct {
	Book    _entity.Book     `json:"book"`
	Reviews []_entity.Review `json:"reviews"`
}

type CreateReviewResponse struct {
	Review _entity.SimplifiedReview `json:"review"`
}

type UpdateReviewRequest struct {
	Star    uint   `json:"star"`
	Content string `json:"content"`
	IsRead  bool   `json:"is_read"`
}

type UpdateReviewResponse struct {
	Review _entity.SimplifiedReview `json:"review"`
}

type GetAllRequestResponse struct {
	Requests []_entity.Request `json:"requests"`
}

type GetAllRequestByUserIdResponse struct {
	Requests []_entity.SimplifiedRequest `json:"requests"`
}

type GetRequestByIdResponse struct {
	Request _entity.Request `json:"request"`
}

type CreateRequestRequest struct {
	BookId uint `json:"book_id"`
}

type CreateRequestResponse struct {
	Request _entity.Request `json:"request"`
}

type UpdateRequestRequest struct {
	ActionCode uint `json:"action_code"`
}

type UpdateRequestResponse struct {
	Request _entity.Request `json:"request"`
}
