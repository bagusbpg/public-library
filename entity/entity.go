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
	Id          uint      `json:"id"`
	Title       string    `json:"title"`
	Author      []Author  `json:"author"`
	Publisher   string    `json:"publisher"`
	Language    string    `json:"language"`
	Pages       uint      `json:"pages"`
	Category    string    `json:"category"`
	ISBN13      string    `json:"isbn13"`
	Description string    `json:"description"`
	Quantity    uint      `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BookItem struct {
	Id     uint
	Book   Book
	Status BookStatus
}

type Author struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type BookStatus struct {
	Id          uint
	Description string
}

type Review struct {
	Id        uint
	User      User
	Book      Book
	Star      uint      `json:"star"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Rent struct {
	Id          uint
	Book        BookItem
	User        User
	Description string `json:"description"`
	Status      RentStatus
	Activity    Activity
	StartAt     time.Time `json:"start_at"`
	ReturnAt    time.Time `json:"return_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RentStatus struct {
	Id          uint
	Description string `json:"status"`
}

type Activity struct {
	Id          uint
	Description string `json:"activity"`
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

type AllWish struct {
	Id        uint      `json:"id"`
	User      User      `json:"user"`
	Title     string    `json:"title"`
	Author    []Author  `json:"author"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
