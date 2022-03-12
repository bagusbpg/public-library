package book

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	_entity "plain-go/public-library/entity"
)

type BookRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (br *BookRepository) GetAllAuthors() (authors []_entity.Author, code int, err error) {
	// prepare statment before execution
	stmt, err := br.db.Prepare(`
		SELECT id, name
		FROM authors
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statment
	row, err := stmt.Query()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer row.Close()

	for row.Next() {
		author := _entity.Author{}

		if err = row.Scan(&author.Id, &author.Name); err != nil {
			log.Println(err)
			code, err = http.StatusInternalServerError, errors.New("internal server error")
			return
		}

		authors = append(authors, author)
	}

	return
}

func (br *BookRepository) CreateNewBook(newBook _entity.Book) (book _entity.Book, code int, err error) {
	// prepare statment before execution
	stmt, err := br.db.Prepare(`
		INSERT INTO books (title, publisher, language, pages, category, isbn13, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	now := time.Now()
	res, err := stmt.Exec(newBook.Title, newBook.Publisher, newBook.Language, newBook.Pages, newBook.Category, newBook.ISBN13, now, now)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	// get new book id
	id, err := res.LastInsertId()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	book = newBook
	book.Id = int(id)
	book.CreatedAt = now
	book.UpdatedAt = now

	return
}

func (br *BookRepository) CreateNewAuthor(newAuthor _entity.Author) (author _entity.Author, code int, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		INSERT INTO authors (name)
		VALUES (?)
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	res, err := stmt.Exec(newAuthor.Name)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	// get new author id
	id, err := res.LastInsertId()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	author = newAuthor
	author.Id = int(id)

	return
}

func (br *BookRepository) CreateBookAuthorJunction(book _entity.Book, author _entity.Author) (code int, err error) {
	// prepare statment before execution
	stmt, err := br.db.Prepare(`
		INSERT INTO book_author_junction (book_id, author_id)
		VALUES (?, ?)
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(book.Id, author.Id)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	return
}

func (br *BookRepository) GetAllBooks() (books []_entity.Book, code int, err error) {
	// prepare statment before execution
	stmt, err := br.db.Prepare(`
		SELECT id, title, publisher, language, pages, category, isbn13, description, created_at, updated_at
		FROM books
		WHERE deleted_at IS NULL
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	for row.Next() {
		book := _entity.Book{}

		if err = row.Scan(&book.Id, &book.Title, &book.Publisher, &book.Language, &book.Pages, &book.Category, &book.ISBN13, &book.Description, &book.CreatedAt, &book.UpdatedAt); err != nil {
			log.Println(err)
			code, err = http.StatusInternalServerError, errors.New("internal server error")
			return
		}

		authors, _, _ := br.getBookAuthors(book.Id)

		book.Author = authors

		books = append(books, book)
	}

	return
}

func (br *BookRepository) getBookAuthors(bookId int) (authors []_entity.Author, code int, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT a.id, a.name
		FROM authors a
		JOIN book_author_juction ba
		ON a.id = ba.author_id
		WHERE ba.book_id = ?
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(bookId)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer row.Close()

	for row.Next() {
		author := _entity.Author{}

		if err = row.Scan(&author.Id, &author.Name); err != nil {
			log.Println(err)
			code, err = http.StatusInternalServerError, errors.New("internal server error")
			return
		}

		authors = append(authors, author)
	}

	return
}
