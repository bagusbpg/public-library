package book

import (
	"database/sql"
	"log"
	"strings"

	_entity "plain-go/public-library/entity"
)

type BookRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (br *BookRepository) GetBookByTitle(title string) (book _entity.Book, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT id, title, publisher, language, pages, category, isbn13, description, created_at, updated_at
		FROM books
		WHERE deleted_at IS NULL
		  AND UPPER(title) LIKE ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(strings.ToUpper(title))

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&book.Id, &book.Title, &book.Publisher, &book.Language, &book.Pages, &book.Category, &book.ISBN13, &book.Description, &book.CreatedAt, &book.UpdatedAt); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (br *BookRepository) GetAuthorByName(name string) (author _entity.Author, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT id, name
		FROM authors
		WHERE UPPER(name) LIKE ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(strings.ToUpper(name))

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&author.Id, &author.Name); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (br *BookRepository) GetAllAuthors() (authors []_entity.Author, err error) {
	// prepare statement
	stmt, err := br.db.Prepare(`
		SELECT id, name
		FROM authors
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statment
	row, err := stmt.Query()

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	for row.Next() {
		author := _entity.Author{}

		if err = row.Scan(&author.Id, &author.Name); err != nil {
			log.Println(err)
			return
		}

		authors = append(authors, author)
	}

	return
}

func (br *BookRepository) CreateNewBook(newBook _entity.Book) (book _entity.Book, err error) {
	// prepare statement
	stmt, err := br.db.Prepare(`
		INSERT INTO books (title, publisher, language, pages, category, isbn13, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	res, err := stmt.Exec(newBook.Title, newBook.Publisher, newBook.Language, newBook.Pages, newBook.Category, newBook.ISBN13, newBook.Description, newBook.CreatedAt, newBook.UpdatedAt)

	if err != nil {
		log.Println(err)
		return
	}

	// get new book id
	id, err := res.LastInsertId()

	if err != nil {
		log.Println(err)
		return
	}

	book = newBook
	book.Id = uint(id)

	return
}

func (br *BookRepository) CreateNewAuthor(newAuthor _entity.Author) (author _entity.Author, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		INSERT INTO authors (name)
		VALUES (?)
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	res, err := stmt.Exec(newAuthor.Name)

	if err != nil {
		log.Println(err)
		return
	}

	// get new author id
	id, err := res.LastInsertId()

	if err != nil {
		log.Println(err)
		return
	}

	author = newAuthor
	author.Id = uint(id)

	return
}

func (br *BookRepository) CreateBookAuthorJunction(book _entity.Book, author _entity.Author) (err error) {
	// prepare statement
	stmt, err := br.db.Prepare(`
		INSERT INTO book_author_junction (book_id, author_id)
		VALUES (?, ?)
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(book.Id, author.Id)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) CreateBookItem(book _entity.Book) (err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		INSERT INTO book_items (book_id, status)
		VALUES (?, 'AVAILABLE')
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(book.Id)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) GetAllBooks() (books []_entity.Book, err error) {
	// prepare statement
	stmt, err := br.db.Prepare(`
		SELECT id, title, publisher, language, pages, category, isbn13, description, created_at, updated_at
		FROM books
		WHERE deleted_at IS NULL
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query()

	if err != nil {
		log.Println(err)
		return
	}

	for row.Next() {
		book := _entity.Book{}

		if err = row.Scan(&book.Id, &book.Title, &book.Publisher, &book.Language, &book.Pages, &book.Category, &book.ISBN13, &book.Description, &book.CreatedAt, &book.UpdatedAt); err != nil {
			log.Println(err)
			return
		}

		books = append(books, book)
	}

	return
}

func (br *BookRepository) GetBookAuthors(bookId uint) (authors []_entity.Author, err error) {
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
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(bookId)

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	for row.Next() {
		author := _entity.Author{}

		if err = row.Scan(&author.Id, &author.Name); err != nil {
			log.Println(err)
			return
		}

		authors = append(authors, author)
	}

	return
}

func (br *BookRepository) GetBookById(bookId uint) (book _entity.Book, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT id, title, publisher, language, pages, category, isbn13, description, created_at, updated_at
		FROM books
		WHERE deleted_at IS NULL
		  AND id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(bookId)

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&book.Id, &book.Title, &book.Publisher, &book.Language, &book.Pages, &book.Category, &book.ISBN13, &book.Description, &book.CreatedAt, &book.UpdatedAt); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (br *BookRepository) CountBookById(bookId uint) (count uint, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT COUNT(id)
		FROM book_items
		WHERE book_id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(bookId)

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&count); err != nil {
			log.Println(err)
			return
		}
	}

	return
}
