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
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Author    []Author
	Publisher Publisher `json:"publisher"`
	Language  Language  `json:"language"`
	Pages     int       `json:"pages"`
	Category  Category
	ISBN13    string    `json:"isbn13"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BookItem struct {
	Id   int
	Book Book
}

type Author struct {
	Id   int
	Name string `json:"author"`
}

type Publisher struct {
	Id   int
	Name string `json:"publisher"`
}

type Language struct {
	Id          int
	Description string `json:"language"`
}

type Category struct {
	Id          int
	Description string `json:"category"`
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
	Book        Book
	User        User
	Description string `json:"description"`
	Status      Status
	Activity    Activity
	StartAt     time.Time `json:"start_at"`
	ReturnAt    time.Time `json:"return_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Status struct {
	Id          int
	Description string `json:"status"`
}

type Activity struct {
	Id          int
	Description string `json:"activity"`
}
