package book

import (
	"log"
	"net/http"
	_bookRepository "plain-go/public-library/datastore/book"
	_entity "plain-go/public-library/entity"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
	"strings"
	"time"
)

type BookUseCase struct {
	repository _bookRepository.Book
}

func New(book _bookRepository.Book) *BookUseCase {
	return &BookUseCase{repository: book}
}

func (buc BookUseCase) CreateBook(req _model.CreateBookRequest) (res _model.CreateBookResponse, code int, message string) {
	// prepare input string
	title := strings.Title(strings.TrimSpace(req.Title))
	publisher := strings.TrimSpace(req.Publisher)
	language := strings.TrimSpace(req.Language)
	category := strings.TrimSpace(req.Category)
	isbn13 := strings.TrimSpace(req.ISBN13)
	description := strings.TrimSpace(req.Description)

	check := []string{title, publisher, language, category, isbn13, description}

	for _, s := range check {
		// check if required input is empty
		if s == "" {
			log.Println("empty input")
			code, message = http.StatusBadRequest, "empty input"
			return
		}

		// check if there is any forbidden character in required field
		if strings.Contains(strings.ReplaceAll(s, " ", ""), ";--") {
			log.Println("forbidden character")
			code, message = http.StatusBadRequest, "forbidden character"
			return
		}
	}

	// check authors
	if len(req.Author) == 0 {
		log.Println("empty input")
		code, message = http.StatusBadRequest, "empty input"
		return
	}

	for _, _author := range req.Author {
		_author.Name = strings.TrimSpace(_author.Name)

		// check if required input is empty
		if _author.Name == "" {
			log.Println("empty input")
			code, message = http.StatusBadRequest, "empty input"
			return
		}

		// check if there is any forbidden character in required field
		if strings.Contains(strings.ReplaceAll(_author.Name, " ", ""), ";--") {
			log.Println("forbidden character")
			code, message = http.StatusBadRequest, "forbidden character"
			return
		}
	}

	// check if pages invalid
	if req.Pages <= 0 {
		log.Println("invalid number of pages")
		code, message = http.StatusBadRequest, "invalid number of pages"
		return
	}

	// check if quantity invalid
	if req.Quantity <= 0 {
		log.Println("invalid number of quantity")
		code, message = http.StatusBadRequest, "invalid number of quantity"
		return
	}

	// check if book is already exist
	newBook, err := buc.repository.GetBookByTitle(title)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	if newBook.Title != "" {
		log.Println("book already exist")
		code, message = http.StatusConflict, "book already exist"
		return
	}

	// prepare input to repository
	now := time.Now()
	newBook.Title = title
	newBook.Publisher = publisher
	newBook.Language = language
	newBook.Pages = req.Pages
	newBook.Category = category
	newBook.ISBN13 = isbn13
	newBook.Description = description
	newBook.CreatedAt = now
	newBook.UpdatedAt = now

	// calling repository
	res.Book, err = buc.repository.CreateNewBook(newBook)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// create author
	for _, _author := range req.Author {
		// calling repository
		author, err := buc.repository.GetAuthorByName(_author.Name)

		// detect failure in repository
		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// if author does not exist, create new
		if author.Name == "" {
			author.Name = _author.Name

			// calling repository
			author, err = buc.repository.CreateNewAuthor(author)

			// detect failure in repository
			if err != nil {
				code, message = http.StatusInternalServerError, "internal server error"
				return
			}
		}

		// if author exist or after author created, create book author junction
		if err = buc.repository.CreateBookAuthorJunction(res.Book, author); err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		res.Book.Author = append(res.Book.Author, author)
	}

	for i := uint(0); i < req.Quantity; i++ {
		if err = buc.repository.CreateBookItem(res.Book); err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}
	}

	// formatting response
	res.Book.Quantity = req.Quantity
	res.Book.CreatedAt, _ = _helper.TimeFormatter(res.Book.CreatedAt)
	res.Book.UpdatedAt, _ = _helper.TimeFormatter(res.Book.UpdatedAt)
	code, message = http.StatusCreated, "success create book"

	return
}

func (buc BookUseCase) GetAllBooks() (res _model.GetAllBooksResponse, code int, message string) {
	// calling repository
	books, err := buc.repository.GetAllBooks()

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, book := range books {
		authors, err := buc.repository.GetBookAuthors(book.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		book.Author = authors

		quantity, err := buc.repository.CountBookById(book.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		book.Quantity = quantity

		res.Books = append(res.Books, book)
	}

	res.Count = uint(len(books))
	code, message = http.StatusOK, "success get all books"

	return
}

func (buc BookUseCase) GetBookById(bookId uint) (res _model.GetBookByIdResponse, code int, message string) {
	// calling repository
	book, err := buc.repository.GetBookById(bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	res.Book = book

	// check if book does not exist
	if res.Book.Title == "" {
		log.Println("book not found")
		code, message = http.StatusNotFound, "book not found"
		return
	}

	// get book count
	res.Book.Quantity, err = buc.repository.CountBookById(bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// get book author
	res.Book.Author, err = buc.repository.GetBookAuthors(bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// formatting response
	res.Book.CreatedAt, _ = _helper.TimeFormatter(res.Book.CreatedAt)
	res.Book.UpdatedAt, _ = _helper.TimeFormatter(res.Book.UpdatedAt)
	code, message = http.StatusOK, "success get book"

	return
}

func (buc BookUseCase) UpdateBook(req _model.UpdateBookRequest, book _entity.Book) (res _model.UpdateBookResponse, code int, message string) {
	// prepare input string
	title := strings.Title(strings.TrimSpace(req.Title))
	publisher := strings.TrimSpace(req.Publisher)
	language := strings.TrimSpace(req.Language)
	category := strings.TrimSpace(req.Category)
	isbn13 := strings.TrimSpace(req.ISBN13)
	description := strings.TrimSpace(req.Description)

	check := []string{title, publisher, language, category, isbn13, description}
	flag := true

	for _, s := range check {
		// check if there is any forbidden character
		if strings.Contains(strings.ReplaceAll(s, " ", ""), ";--") {
			log.Println("forbidden character")
			code, message = http.StatusBadRequest, "forbidden chacarter"
			return
		}

		// check if any field is updated
		if s != "" {
			flag = false
		}
	}

	if title != "" {
		book.Title = title
	}

	if publisher != "" {
		book.Publisher = publisher
	}

	if language != "" {
		book.Language = language
	}

	if category != "" {
		book.Category = category
	}

	if isbn13 != "" {
		book.ISBN13 = isbn13
	}

	if description != "" {
		book.Description = description
	}

	// if authors are updated, clear existing ones first
	if len(req.Author) > 0 {
		book.Author = []_entity.Author{}
	}

	for _, _author := range req.Author {
		author := _entity.Author{}
		author.Name = strings.TrimSpace(_author.Name)

		// check if there is any forbidden character in required field
		if strings.Contains(strings.ReplaceAll(_author.Name, " ", ""), ";--") {
			log.Println("forbidden character")
			code, message = http.StatusBadRequest, "forbidden character"
			return
		}

		book.Author = append(book.Author, author)
		flag = false
	}

	// check if pages invalid
	if req.Pages < 0 {
		log.Println("invalid number of pages")
		code, message = http.StatusBadRequest, "invalid number of pages"
		return
	} else if req.Pages > 0 {
		book.Pages = req.Pages
		flag = false
	}

	// check if no field is updated
	if flag {
		log.Println("no update was performed")
		code, message = http.StatusBadRequest, "no update was performed"
		return
	}

	// calling repository
	book.UpdatedAt = time.Now()
	_book, err := buc.repository.UpdateBook(book)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// formatting response
	res.Book = _book
	res.Book.Id = book.Id
	res.Book.UpdatedAt, _ = _helper.TimeFormatter(res.Book.UpdatedAt)
	code, message = http.StatusOK, "success update book"

	return
}
