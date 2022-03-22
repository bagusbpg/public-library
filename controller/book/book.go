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
		res, code, message := bc.usecase.GetAllBooks()

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
