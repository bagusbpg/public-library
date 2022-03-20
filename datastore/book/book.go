package book

import (
	"database/sql"
	"log"
	"strings"
	"time"

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
		JOIN book_author_junction ba
		ON a.id = ba.author_id
		WHERE ba.book_id = ?
		  AND ba.deleted_at IS NULL
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

func (br *BookRepository) UpdateBook(updatedBook _entity.Book) (book _entity.Book, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		UPDATE books
		SET title = ?, publisher = ?, language = ?, pages = ?, category = ?, isbn13 = ?, description = ?, updated_at = ?
		WHERE id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(updatedBook.Title, updatedBook.Publisher, updatedBook.Language, updatedBook.Pages, updatedBook.Category, updatedBook.ISBN13, updatedBook.Description, updatedBook.UpdatedAt, updatedBook.Id)

	if err != nil {
		log.Println(err)
		return
	}

	book = updatedBook

	return
}

func (br *BookRepository) DeleteBookAuthorJunction(book _entity.Book, author _entity.Author) (err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		UPDATE book_author_junction
		SET deleted_at = ?
		WHERE book_id = ?
		  AND author_id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(time.Now(), book.Id, author.Id)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) AddBookToFavorite(userId uint, bookId uint) (err error) {
	// prepare statement
	stmt, err := br.db.Prepare(`
		INSERT INTO favorites (user_id, book_id, created_at)
		VALUES (?, ?, ?)
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(userId, bookId, time.Now())

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) RemoveBookFromFavorite(userId uint, bookId uint) (err error) {
	// prepare statement
	stmt, err := br.db.Prepare(`
		UPDATE favorites
		SET deleted_at = ?
		WHERE user_id = ?
	  	  AND book_id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(time.Now(), userId, bookId)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) GetAllFavorites(userId uint) (favorites []_entity.Favorite, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT id, book_id, created_at
		FROM favorites
		WHERE user_id = ?
		  AND deleted_at IS NULL
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(userId)

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	for row.Next() {
		favorite := _entity.Favorite{}

		if err = row.Scan(&favorite.Id, &favorite.Book.Id, &favorite.CreatedAt); err != nil {
			log.Println(err)
			return
		}

		favorites = append(favorites, favorite)
	}

	return
}

func (br *BookRepository) AddBookToWishlist(userId uint, newWish _entity.Wish) (wish _entity.Wish, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		INSERT INTO wishlists (user_id, title, category, created_at)
		VALUES (?, ?, ?, ?)
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	res, err := stmt.Exec(userId, newWish.Title, newWish.Category, newWish.CreatedAt)

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

	wish = newWish
	wish.Id = uint(id)

	return
}

func (br *BookRepository) RemoveBookFromWishlist(wishId uint) (err error) {
	// prepare statement
	stmt, err := br.db.Prepare(`
		UPDATE wishlists
		SET deleted_at = ?
		WHERE wish_id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(time.Now(), wishId)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) GetAllWishes(userId uint) (wishes []_entity.Wish, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT id, title, category, created_at
		FROM wishlists
		WHERE user_id = ?
		  AND deleted_at IS NULL
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(userId)

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	for row.Next() {
		wish := _entity.Wish{}

		if err = row.Scan(&wish.Id, &wish.Title, &wish.Category, &wish.CreatedAt); err != nil {
			log.Println(err)
			return
		}

		wishes = append(wishes, wish)
	}

	return
}

func (br *BookRepository) CreateWishAuthorJunction(wish _entity.Wish, author _entity.Author) (err error) {
	// prepare statement
	stmt, err := br.db.Prepare(`
		INSERT INTO wish_author_junction (wish_id, author_id)
		VALUES (?, ?)
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(wish.Id, author.Id)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) GetWishAuthors(wishId uint) (authors []_entity.Author, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT a.id, a.name
		FROM authors a
		JOIN wish_author_junction wa
		ON a.id = wa.author_id
		WHERE wa.wish_id = ?
		  AND wa.deleted_at IS NULL
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(wishId)

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
