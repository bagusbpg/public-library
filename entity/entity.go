package entity

import "time"

type User struct {
	Id        int       `json:"id"`
	Role      string    `json:"role"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Book struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Author      []Author
	Publisher   string    `json:"publisher"`
	Language    string    `json:"language"`
	Pages       int       `json:"pages"`
	Category    string    `json:"category"`
	ISBN13      string    `json:"isbn13"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BookItem struct {
	Id     int
	Book   Book
	Status BookStatus
}

type Author struct {
	Id   int
	Name string `json:"author"`
}

type BookStatus struct {
	Id          int
	Description string
}

type Review struct {
	Id        int
	User      User
	Book      Book
	Star      int       `json:"star"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Rent struct {
	Id          int
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
	Id          int
	Description string `json:"status"`
}

type Activity struct {
	Id          int
	Description string `json:"activity"`
}
