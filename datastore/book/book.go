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
		VALUES (?, 'availabe')
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

func (br *BookRepository) DeleteBook(bookId uint) (err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		UPDATE books
		SET deleted_at = ?
		WHERE id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(time.Now(), bookId)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) AddBookToFavorite(userId uint, bookId uint) (favorite _entity.Favorite, err error) {
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
	now := time.Now()
	res, err := stmt.Exec(userId, bookId, now)

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

	favorite.Id = uint(id)
	favorite.CreatedAt = now

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

func (br *BookRepository) GetAllFavoritesByUserId(userId uint) (favorites []_entity.Favorite, err error) {
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

func (br *BookRepository) CountFavoritesByBookId(bookId uint) (count uint, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT COUNT(id)
		FROM favorites
		WHERE book_id = ?
		GROUP BY book_id
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

func (br *BookRepository) GetAllWishes() (wishes []_entity.AllWish, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT id, user_id, title, category, note, created_at, updated_at
		FROM wishlistst
		WHERE deleted_at IS NULL
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	row, err := stmt.Query()

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	for row.Next() {
		wish := _entity.AllWish{}

		if err = row.Scan(&wish.Id, &wish.User.Id, &wish.Title, &wish.Category, &wish.Note, &wish.CreatedAt, &wish.UpdatedAt); err != nil {
			log.Println(err)
			return
		}

		wishes = append(wishes, wish)
	}

	return
}

func (br *BookRepository) AddBookToWishlist(userId uint, newWish _entity.Wish) (wish _entity.Wish, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		INSERT INTO wishlists (user_id, title, category, note, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	res, err := stmt.Exec(userId, newWish.Title, newWish.Category, newWish.Note, newWish.CreatedAt, newWish.UpdatedAt)

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
		WHERE id = ?
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

func (br *BookRepository) GetWishesByUserId(userId uint) (wishes []_entity.Wish, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT id, title, category, note, created_at, updated_at
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

		if err = row.Scan(&wish.Id, &wish.Title, &wish.Category, &wish.Note, &wish.CreatedAt, &wish.UpdatedAt); err != nil {
			log.Println(err)
			return
		}

		wishes = append(wishes, wish)
	}

	return
}

func (br *BookRepository) GetWishById(userId uint, wishId uint) (wish _entity.Wish, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT id, title, category, note, created_at, updated_at
		FROM wishlists
		WHERE deleted_at IS NULL
		  AND id = ?
		  AND user_id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(wishId, userId)

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&wish.Id, &wish.Title, &wish.Category, &wish.Note, &wish.CreatedAt, &wish.UpdatedAt); err != nil {
			log.Println(err)
			return
		}
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

func (br *BookRepository) UpdateWish(updatedWish _entity.Wish) (wish _entity.Wish, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		UPDATE wishlists
		SET title = ?, category = ?, note = ?, updated_at = ?
		WHERE id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(updatedWish.Title, updatedWish.Category, updatedWish.Note, updatedWish.UpdatedAt, updatedWish.Id)

	if err != nil {
		log.Println(err)
		return
	}

	wish = updatedWish

	return
}

func (br *BookRepository) DeleteWishAuthorJunction(wish _entity.Wish, author _entity.Author) (err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		UPDATE wish_author_junction
		SET deleted_at = ?
		WHERE wish_id = ?
		  AND author_id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(time.Now(), wish.Id, author.Id)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) GetAllReviews() (reviews []_entity.AllReview, err error) {
	// prepare statment before execution
	stmt, err := br.db.Prepare(`
		SELECT id, user_id, book_id, star, content, flag, created_at, updated_at
		FROM reviews
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

	defer row.Close()

	for row.Next() {
		review := _entity.AllReview{}

		if err = row.Scan(&review.Id, &review.User.Id, &review.Book.Id, &review.Star, &review.Content, &review.Flag, &review.CreatedAt, &review.UpdatedAt); err != nil {
			log.Println(err)
			return
		}

		reviews = append(reviews, review)
	}

	return
}

func (br *BookRepository) CreateReview(newReview _entity.AllReview) (review _entity.AllReview, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		INSERT INTO reviews (user_id, book_id, star, content, flag, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	flag := 0
	res, err := stmt.Exec(newReview.User.Id, newReview.Book.Id, newReview.Star, newReview.Content, flag, newReview.CreatedAt, newReview.UpdatedAt)

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

	review = newReview
	review.Id = uint(id)

	return
}

func (br *BookRepository) GetReviewByReviewId(reviewId uint) (review _entity.AllReview, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT id, user_id, book_id, star, content, flag, created_at, updated_at
		FROM reviews
		WHERE id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(reviewId)

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&review.Id, &review.User.Id, &review.Book.Id, &review.Star, &review.Content, &review.Flag, &review.CreatedAt, &review.UpdatedAt); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (br *BookRepository) GetAllReviewsByBookId(bookId uint) (reviews []_entity.Review, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT id, user_id, star, content, created_at, updated_at
		FROM reviews
		WHERE book_id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	row, err := stmt.Query(bookId)

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	for row.Next() {
		review := _entity.Review{}

		if err = row.Scan(&review.Id, &review.User.Id, &review.Star, &review.Content, &review.CreatedAt, &review.UpdatedAt); err != nil {
			log.Println(err)
			return
		}

		reviews = append(reviews, review)
	}

	return
}

func (br *BookRepository) UpdateReview(updatedReview _entity.AllReview) (review _entity.AllReview, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		UPDATE reviews
		SET star = ?, content = ?, flag = ?, updated_at = ?
		WHERE id = ?
		  AND deleted_at IS NULL
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(updatedReview.Star, updatedReview.Content, &updatedReview.Flag, updatedReview.UpdatedAt, updatedReview.Id)

	if err != nil {
		log.Println(err)
		return
	}

	review = updatedReview

	return
}

func (br *BookRepository) UpdateReviewStatus(flag uint, reviewId uint) (err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		UPDATE reviews
		SET flag = ?
		WHERE id = ?
	  	  AND deleted_at IS NULL
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(flag, reviewId)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) DeleteReview(reviewId uint) (err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		UPDATE reviews
		SET deleted_at = ?
		WHERE id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(reviewId)

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (br *BookRepository) CountStarsByBookId(bookId uint) (averageStar float64, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT AVG(star)
		FROM reviews
		WHERE book_id = ?
		GROUP BY book_id
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
		if err = row.Scan(&averageStar); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (br *BookRepository) GetBookByItemId(itemId uint) (book _entity.Book, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT b.id, b.title, b.publisher, b.language, b.pages, b.category, b.isbn13, b.description, b.created_at, b.updated_at
		FROM books b
		JOIN book_items bi
		ON b.id = bi.book_id
		WHERE bi.id = ?
		  AND b.deleted_at IS NULL
		LIMIT 1
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(itemId)

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

func (br *BookRepository) GetAvailableBookByBookId(bookId uint) (bookItemId uint, err error) {
	// prepare statement before execution
	stmt, err := br.db.Prepare(`
		SELECT bi.id
		FROM book_items bi
		JOIN books b
		ON bi.book_id = b.id
		WHERE b.id = ?
		  AND b.deleted_at IS NULL
		  AND bi.status = 'available'
		ORDER BY bi.id ASC
		LIMIT 1
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

	if row.Next() {
		if err = row.Scan(&bookItemId); err != nil {
			log.Println(err)
			return
		}
	}

	return
}
