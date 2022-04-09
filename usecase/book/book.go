package book

import (
	"log"
	"net/http"
	"net/url"
	_bookRepository "plain-go/public-library/datastore/book"
	_entity "plain-go/public-library/entity"
	_helper "plain-go/public-library/helper"
	_model "plain-go/public-library/model"
	"strconv"
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

func (buc BookUseCase) GetAllBooks(query url.Values) (res _model.GetAllBooksResponse, code int, message string) {
	// default parameters
	params := _model.GetAllBooksRequest{}
	params.Page = 1
	params.Records = 9
	params.Category = "*"
	params.Keyword = "*"
	params.SortBy = "*"
	params.SortMode = "random"

	if value, exist := query["page"]; exist {
		page, err := strconv.Atoi(value[0])

		if err != nil {
			log.Println(err)
			code, message = http.StatusBadRequest, "invalid page"
			return
		}

		if page < 1 {
			log.Println("invalid page")
			code, message = http.StatusBadRequest, "invalid page"
			return
		}

		params.Page = page
	}

	mapRecords := map[int]interface{}{9: nil, 15: nil, 30: nil, 60: nil, 90: nil}

	if value, exist := query["records"]; exist {
		records, err := strconv.Atoi(value[0])

		if err != nil {
			log.Println(err)
			code, message = http.StatusBadRequest, "invalid number of records"
			return
		}

		if _, exist := mapRecords[records]; !exist {
			log.Println("unaccepted number of records")
			code, message = http.StatusBadRequest, "unaccepted number of records"
			return
		}

		params.Records = records
	}

	if value, exist := query["category"]; exist {
		params.Category = value[0]
	}

	if value, exist := query["keyword"]; exist {
		params.Keyword = strings.Join(value, " ")
	}

	mapSort := map[string]interface{}{"star": nil, "review": nil, "read": nil}

	if value, exist := query["sort"]; exist {
		if _, exist := mapSort[value[0]]; !exist {
			log.Println("unaccepted sorting criteria")
			code, message = http.StatusBadRequest, "unaccepted sorting criteria"
			return
		}

		params.SortBy = value[0]
	}

	mapMode := map[string]interface{}{"asc": nil, "desc": nil}

	if value, exist := query["mode"]; exist {
		if _, exist := mapMode[value[0]]; !exist {
			log.Println("unaccepted sorting mode")
			code, message = http.StatusBadRequest, "unaccepted sorting mode"
			return
		}

		params.SortMode = value[0]
	}

	// calling repository
	books, err := buc.repository.GetAllBooks(params)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	for _, book := range books {
		// get book author
		book.Author, err = buc.repository.GetBookAuthors(book.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// get book quantity
		book.Quantity, err = buc.repository.CountBookById(book.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		// get book favorite count
		book.FavoriteCount, err = buc.repository.CountFavoritesByBookId(book.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		averageStar, err := buc.repository.CountStarsByBookId(book.Id)

		if err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		if averageStar == 0 {
			book.AverageStar = nil
		} else {
			book.AverageStar = averageStar
		}

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

	// get favorite count
	res.Book.FavoriteCount, err = buc.repository.CountFavoritesByBookId(bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// get average star
	averageStar, err := buc.repository.CountStarsByBookId(book.Id)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	book.AverageStar = _helper.NilHandler(averageStar)

	// formatting response
	res.Book.CreatedAt, _ = _helper.TimeFormatter(res.Book.CreatedAt)
	res.Book.UpdatedAt, _ = _helper.TimeFormatter(res.Book.UpdatedAt)
	code, message = http.StatusOK, "success get book"

	return
}

func (buc BookUseCase) UpdateBook(req _model.UpdateBookRequest, bookId uint) (res _model.UpdateBookResponse, code int, message string) {
	// check book existence
	book, err := buc.repository.GetBookById(bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// check if book does not exist
	if book.Title == "" {
		log.Println("book not found")
		code, message = http.StatusNotFound, "book not found"
		return
	}

	// get book's authors
	book.Author, err = buc.repository.GetBookAuthors(bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

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
	}

	if title != "" && title != book.Title {
		book.Title = title
		flag = false
	}

	if publisher != "" && publisher != book.Publisher {
		book.Publisher = publisher
		flag = false
	}

	if language != "" && language != book.Language {
		book.Language = language
		flag = false
	}

	if category != "" && category != book.Category {
		book.Category = category
		flag = false
	}

	if isbn13 != "" && isbn13 != book.ISBN13 {
		book.ISBN13 = isbn13
		flag = false
	}

	if description != "" && description != book.Description {
		book.Description = description
		flag = false
	}

	// if authors are updated
	existingAuthors := map[string]interface{}{}
	updatedAuthors := map[string]interface{}{}

	if len(req.Author) > 0 {
		for _, author := range book.Author {
			existingAuthors[author.Name] = nil
		}
	}

	for _, _author := range req.Author {
		author := _entity.Author{}
		author.Name = strings.TrimSpace(_author.Name)

		// check if there is any forbidden character in required field
		if strings.Contains(strings.ReplaceAll(author.Name, " ", ""), ";--") {
			log.Println("forbidden character")
			code, message = http.StatusBadRequest, "forbidden character"
			return
		}

		updatedAuthors[author.Name] = nil
	}

	createdAuthors := []_entity.Author{}
	droppedAuthors := []_entity.Author{}

	for name := range updatedAuthors {
		if _, exist := existingAuthors[name]; !exist {
			createdAuthors = append(createdAuthors, _entity.Author{Name: name})
		}
	}

	for name := range existingAuthors {
		if _, exist := updatedAuthors[name]; !exist {
			droppedAuthors = append(droppedAuthors, _entity.Author{Name: name})
		}
	}

	if len(createdAuthors) != 0 || len(droppedAuthors) != 0 {
		flag = false
	}

	if req.Pages > 0 && req.Pages != book.Pages {
		book.Pages = req.Pages
		flag = false
	}

	// check if no field is updated
	if flag {
		log.Println("no update was performed")
		code, message = http.StatusBadRequest, "no update was performed"
		return
	}

	// create author if len(createdAuthors) > 0
	for _, _author := range createdAuthors {
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
		if err = buc.repository.CreateBookAuthorJunction(book, author); err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		book.Author = append(book.Author, author)
	}

	// drop author if len(droppedAuthors) > 0
	for _, _author := range droppedAuthors {
		_author, _ = buc.repository.GetAuthorByName(_author.Name)

		if err := buc.repository.DeleteBookAuthorJunction(book, _author); err != nil {
			code, message = http.StatusInternalServerError, "internal server error"
			return
		}

		book.Author = _helper.RemoveAuthor(book.Author, _author)
	}

	// calling repository
	book.UpdatedAt = time.Now()
	res.Book, err = buc.repository.UpdateBook(book)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// formatting response
	res.Book.Id = book.Id
	res.Book.FavoriteCount, err = buc.repository.CountFavoritesByBookId(bookId)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	averageStar, err := buc.repository.CountStarsByBookId(book.Id)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	res.Book.AverageStar = _helper.NilHandler(averageStar)
	res.Book.CreatedAt, _ = _helper.TimeFormatter(res.Book.CreatedAt)
	res.Book.UpdatedAt, _ = _helper.TimeFormatter(res.Book.UpdatedAt)
	code, message = http.StatusOK, "success update book"

	return
}

func (buc BookUseCase) DeleteBook(bookId uint) (code int, message string) {
	// check book existence
	book, err := buc.repository.GetBookById(bookId)

	// detect failure in repository
	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	// check if book does not exist
	if book.Title == "" {
		log.Println("book not found")
		code, message = http.StatusNotFound, "book not found"
		return
	}

	// calling repository
	err = buc.repository.DeleteBook(bookId)

	if err != nil {
		code, message = http.StatusInternalServerError, "internal server error"
		return
	}

	code, message = http.StatusOK, "success delete book"

	return
}
