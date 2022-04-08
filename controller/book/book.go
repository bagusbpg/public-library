package book

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_mw "plain-go/public-library/app/middleware"
	_model "plain-go/public-library/model"
	_bookUseCase "plain-go/public-library/usecase/book"
	"strconv"
)

type BookController struct {
	usecase _bookUseCase.Book
}

func New(book _bookUseCase.Book) *BookController {
	return &BookController{usecase: book}
}

func (bc BookController) Create() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "internal server error", nil)
			return
		}

		defer r.Body.Close()

		req := _model.CreateBookRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := bc.usecase.CreateBook(req)

		if code != http.StatusCreated {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (bc BookController) GetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		mapRecords := map[int]interface{}{9: nil, 15: nil, 30: nil, 60: nil, 90: nil}
		query := r.URL.Query()

		params := _model.GetAllBooksRequest{}
		params.Page = 1
		params.Records = 9

		if value, exist := query["page"]; exist {
			page, err := strconv.Atoi(value[0])

			if err != nil {
				log.Println(err)
				code, message := http.StatusBadRequest, "invalid page"
				_model.CreateResponse(rw, code, message, nil)
				return
			}

			if page < 1 {
				log.Println("invalid page")
				code, message := http.StatusBadRequest, "invalid page"
				_model.CreateResponse(rw, code, message, nil)
				return
			}

			params.Page = page
		}

		if value, exist := query["records"]; exist {
			records, err := strconv.Atoi(value[0])

			if err != nil {
				log.Println(err)
				code, message := http.StatusBadRequest, "invalid number of records"
				_model.CreateResponse(rw, code, message, nil)
				return
			}

			if _, exist := mapRecords[records]; !exist {
				log.Println("unaccepted number of records")
				code, message := http.StatusBadRequest, "unaccepted number of records"
				_model.CreateResponse(rw, code, message, nil)
				return
			}

			params.Records = records
		}

		res, code, message := bc.usecase.GetAllBooks(params)

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (bc BookController) Get() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bookId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		existing, code, message := bc.usecase.GetBookById(uint(bookId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, existing.Book)
	}
}

func (bc BookController) Update() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bookId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusInternalServerError, "failed to read request body", nil)
			return
		}

		defer r.Body.Close()

		req := _model.UpdateBookRequest{}

		if err = json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			_model.CreateResponse(rw, http.StatusBadRequest, "failed to bind request body", nil)
			return
		}

		res, code, message := bc.usecase.UpdateBook(req, uint(bookId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, res)
	}
}

func (bc BookController) Delete() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bookId, _ := strconv.Atoi(_mw.GetParam(r)[0])

		code, message := bc.usecase.DeleteBook(uint(bookId))

		_model.CreateResponse(rw, code, message, nil)
	}
}
