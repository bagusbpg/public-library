package entity

import "time"

type User struct {
	Id        uint      `json:"id"`
	Role      string    `json:"role"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Book struct {
	Id            uint        `json:"id"`
	Title         string      `json:"title"`
	Author        []Author    `json:"author"`
	Publisher     string      `json:"publisher"`
	Language      string      `json:"language"`
	Pages         uint        `json:"pages"`
	Category      string      `json:"category"`
	ISBN13        string      `json:"isbn13"`
	Description   string      `json:"description"`
	Quantity      uint        `json:"quantity"`
	FavoriteCount uint        `json:"favorite_count"`
	AverageStar   interface{} `json:"average_star"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	// ReadCount     uint        `json:"read_count"`
	// Available     uint        `json:"available"`
}

type BookItem struct {
	Id     int    `json:"book_item_id"`
	Book   Book   `json:"book_item_detail"`
	Status string `json:"book_item_status"`
}

type Author struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type Review struct {
	Id        uint      `json:"id"`
	User      User      `json:"reviewer"`
	Star      uint      `json:"star"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SimplifiedReview struct {
	Id        uint      `json:"id"`
	User      User      `json:"reviewer"`
	Book      Book      `json:"book_reviewed"`
	Star      uint      `json:"star"`
	Content   string    `json:"content"`
	Flag      uint      `json:"flag"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Request struct {
	Id        uint          `json:"id"`
	BookItem  BookItem      `json:"book_item"`
	User      User          `json:"user"`
	Status    RequestStatus `json:"status"`
	Extended  uint          `json:"extended"`
	CreatedAt time.Time     `json:"created_at"`
	StartAt   interface{}   `json:"start_at"`
	FinishAt  interface{}   `json:"finish_at"`
	ReturnAt  interface{}   `json:"return_at"`
	CancelAt  interface{}   `json:"cancel_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type SimplifiedRequest struct {
	Id        uint          `json:"id"`
	BookItem  BookItem      `json:"book_item"`
	Status    RequestStatus `json:"status"`
	Extended  uint          `json:"extended"`
	CreatedAt time.Time     `json:"created_at"`
	StartAt   interface{}   `json:"start_at"`
	FinishAt  interface{}   `json:"finish_at"`
	ReturnAt  interface{}   `json:"return_at"`
	CancelAt  interface{}   `json:"cancel_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type RequestStatus struct {
	Id          uint
	Description string `json:"status"`
}

type Favorite struct {
	Id        uint      `json:"id"`
	Book      Book      `json:"book"`
	CreatedAt time.Time `json:"created_at"`
}

type Wish struct {
	Id        uint      `json:"id"`
	Title     string    `json:"title"`
	Author    []Author  `json:"author"`
	Category  string    `json:"category"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SimplifiedWish struct {
	Id        uint      `json:"id"`
	User      User      `json:"user"`
	Title     string    `json:"title"`
	Author    []Author  `json:"author"`
	Category  string    `json:"category"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
