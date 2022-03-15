package book

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	_model "plain-go/public-library/model"
	_bookUseCase "plain-go/public-library/usecase/book"
)

type BookController struct {
	usecase _bookUseCase.Book
}

func New(book _bookUseCase.Book) *BookController {
	return &BookController{usecase: book}
}

func (bc BookController) Create() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("method not allowed")
			_model.CreateResponse(rw, http.StatusMethodNotAllowed, "method not allowed", nil)
			return
		}

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
	}
}
