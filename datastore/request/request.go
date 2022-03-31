package request

import (
	"database/sql"
	"log"
	_entity "plain-go/public-library/entity"
	"strings"
)

type RequestRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *RequestRepository {
	return &RequestRepository{db: db}
}

func (rr RequestRepository) GetRequestByUserIdAndBookId(userId uint, bookId uint) (requests []_entity.Request, err error) {
	// prepare statement before execution
	stmt, err := rr.db.Prepare(`
		SELECT r.id, r.user_id, r.book_item_id, r.status_id, rs.description, r.created_at, r.start_at, r.return_at, r.updated_at
		FROM requests r
		JOIN book_items bi
		ON r.book_items_id = bi.id
		JOIN books b
		ON bi.book_id = b.id
		JOIN request_status rs
		ON r.status_id = rs.id
		WHERE r.user_id = ?
		  AND b.id = ?
		  AND r.cancel_at IS NULL
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(bookId, userId)

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	for row.Next() {
		request := _entity.Request{}

		if err = row.Scan(&request.Id, &request.User.Id, &request.Book.Id, &request.Status.Id, &request.Status.Description, &request.StartAt, &request.ReturnAt, &request.UpdatedAt); err != nil {
			log.Println(err)
			return
		}

		requests = append(requests, request)
	}

	return
}

func (rr RequestRepository) GetAvailableBookByBookId(bookId uint) (bookItemId uint, err error) {
	// prepare statement before execution
	stmt, err := rr.db.Prepare(`
		SELECT bi.id
		FROM book_items bi
		JOIN books b
		ON bi.book_id = b.id
		WHERE b.id = ?
		  AND b.deleted_at IS NULL
		  AND bi.status = 0
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

func (rr RequestRepository) CreateNewRequest(newRequest _entity.Request) (request _entity.Request, err error) {
	// prepare statement before execution
	stmt, err := rr.db.Prepare(`
		INSERT INTO requests (user_id, book_item_id, status_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	res, err := stmt.Exec(newRequest.User.Id, newRequest.Book.Id, newRequest.Status.Id, newRequest.CreatedAt, newRequest.UpdatedAt)

	if err != nil {
		log.Println(err)
		return
	}

	// get new request id
	id, err := res.LastInsertId()

	if err != nil {
		log.Println(err)
		return
	}

	request = newRequest
	request.Id = uint(id)

	return
}

func (rr RequestRepository) GetRequestByRequestId(requestId uint) (request _entity.Request, err error) {
	// prepare statement before execution
	stmt, err := rr.db.Prepare(`
		SELECT r.id, r.user_id, r.book_id, r.status_id, rs.description, r.created_at, r.start_at, r.return_at, r.cancel_at, r.updated_at
		FROM requests r
		JOIN request_status rs
		ON r.status_id = rs.id
		WHERE id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(requestId)

	if err != nil {
		log.Println(err)
		return
	}

	if row.Next() {
		if err = row.Scan(&request.Id, &request.User.Id, &request.Book.Id, &request.Status.Id, &request.Status.Description, &request.StartAt, &request.ReturnAt, &request.CancelAt, &request.UpdatedAt); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (rr RequestRepository) GetAll() (requests []_entity.Request, err error) {
	// prepare statement before execution
	stmt, err := rr.db.Prepare(`
		SELECT r.id, r.book_item_id, r.user_id, r.status_id, rs.description, r.created_at, r.start_at, r.return_at, r.cancel_at, r.updated_at
		FROM requests r
		JOIN request_status rs
		ON r.status_id = rs.id
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
		request := _entity.Request{}

		if err = row.Scan(&request.Id, &request.Book.Id, &request.User.Id, &request.Status.Id, &request.CreatedAt, &request.StartAt, &request.ReturnAt, &request.CancelAt, &request.UpdatedAt); err != nil {
			log.Println(err)
			return
		}

		requests = append(requests, request)
	}

	return
}

func (rr RequestRepository) GetRequestStatusId(statusDesc string) (statusId uint, err error) {
	// prepare statement before execution
	stmt, err := rr.db.Prepare(`
		SELECT id
		FROM request_status
		WHERE UPPER(description) = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(strings.ToUpper(statusDesc))

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&statusId); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (rr RequestRepository) Update() {

}
