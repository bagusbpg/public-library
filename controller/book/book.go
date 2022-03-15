package book

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_model "plain-go/public-library/model"
	_bookUseCase "plain-go/public-library/usecase/book"
	"strconv"
	"strings"
)

type BookController struct {
	usecase _bookUseCase.Book
}

func New(book _bookUseCase.Book) *BookController {
	return &BookController{usecase: book}
}

func (bc BookController) CreateGetAll() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			res, code, message := bc.usecase.GetAllBooks()

			if code != http.StatusOK {
				_model.CreateResponse(rw, code, message, nil)
				return
			}

			_model.CreateResponse(rw, code, message, res)
		case http.MethodPost:
			body, err := ioutil.ReadAll(r.Body)

			if err != nil {
				log.Println(err)
				_model.CreateResponse(rw, http.StatusInternalServerError, "internal server error", nil)
				return
			}

			defer r.Body.Close()

			if contentType := r.Header.Get("content-type"); contentType != "application/json" {
				log.Println("unsupported content type")
				_model.CreateResponse(rw, http.StatusUnsupportedMediaType, "unsupported content type", nil)
				return
			}

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
		default:
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
		}
	}
}

func (bc BookController) GetUpdateDelete() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bookId, _ := strconv.Atoi(strings.SplitAfter(r.URL.Path, "/")[2])

		existing, code, message := bc.usecase.GetBookById(uint(bookId))

		if code != http.StatusOK {
			_model.CreateResponse(rw, code, message, nil)
			return
		}

		_model.CreateResponse(rw, code, message, existing.Book)
	}
}
