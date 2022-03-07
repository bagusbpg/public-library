package user

import (
	"errors"
	"log"
	"net/http"
	_entity "plain-go/public-library/entity"
	"time"

	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateNewUser(newUser _entity.User) (user _entity.User, code int, err error) {
	// check email existence
	if code, err = ur.checkEmailExistence(newUser.Email); err != nil {
		return
	}

	// check phone existence
	if code, err = ur.checkPhoneExistence(newUser.Phone); err != nil {
		return
	}

	// prepare statement before query or execution
	stmt, err := ur.db.Prepare(`
		INSERT INTO users (role, name, email, phone, password, created_at, updated_at)
		VALUES ('Member', ?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	time := time.Now()
	res, err := stmt.Exec(newUser.Name, newUser.Email, newUser.Phone, newUser.Password, time, time)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	// get new user id
	id, err := res.LastInsertId()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	user = newUser
	user.Id = int(id)
	user.Role = "Member"
	user.CreatedAt = time
	user.UpdatedAt = time

	return
}

func (ur *UserRepository) GetUserDetail(userId int) (user _entity.User, code int, err error) {
	// prepare statement before query or execution
	stmt, err := ur.db.Prepare(`
		SELECT role, name, email, phone, password
		FROM users
		WHERE deleted_at IS NULL
		  AND id = ?
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(userId)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&user.Role, &user.Name, &user.Email, &user.Phone, &user.Password); err != nil {
			log.Println(err)
			code, err = http.StatusInternalServerError, errors.New("internal server error")
			return
		}
	}

	if user == (_entity.User{}) {
		code, err = http.StatusNotFound, errors.New("user not found")
		log.Println(err)
		return
	}

	return
}

func (ur *UserRepository) UpdateUserDetail(updatedUser _entity.User) (user _entity.User, code int, err error) {
	// check email existence
	if code, err = ur.checkEmailExistence(updatedUser.Email); err != nil {
		return
	}

	// check phone existence
	if code, err = ur.checkPhoneExistence(updatedUser.Phone); err != nil {
		return
	}

	// prepare statement before query or execution
	stmt, err := ur.db.Prepare(`
		UPDATE users
		SET name = ?, email = ?, phone = ?, password = ?, updated_at = CURRENT_TIMESTAMP
		WHERE deleted_at IS NULL
		  AND id = ?
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(updatedUser.Name, updatedUser.Email, updatedUser.Phone, updatedUser.Password, updatedUser.Id)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	user = updatedUser

	return
}

func (ur *UserRepository) DeleteUser(userId int) (code int, err error) {
	// prepare statement before query or execution
	stmt, err := ur.db.Prepare(`
		UPDATE users
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE deleted_at IS NULL
		  AND id = ?
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	res, err := stmt.Exec(userId)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	// check user existence
	row, err := res.RowsAffected()

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	if row == 0 {
		code, err = http.StatusNotFound, errors.New("user not found")
		log.Println(err)
		return
	}

	return
}

func (ur UserRepository) checkEmailExistence(email string) (code int, err error) {
	// prepare statement before query or execution
	stmt, err := ur.db.Prepare(`
		SELECT id
		FROM users
		WHERE deleted_at IS NULL
		  AND email = ?

	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(email)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer row.Close()

	id := 0

	if row.Next() {
		if err = row.Scan(&id); err != nil {
			log.Println(err)
			code, err = http.StatusInternalServerError, errors.New("internal server error")
			return
		}
	}

	if id != 0 {
		code, err = http.StatusConflict, errors.New("email already used")
		log.Println(err)
		return
	}

	return
}

func (ur UserRepository) checkPhoneExistence(phone string) (code int, err error) {
	// prepare statement before query or execution
	stmt, err := ur.db.Prepare(`
		SELECT id
		FROM users
		WHERE deleted_at IS NULL
		  AND phone = ?
	`)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(phone)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	defer row.Close()

	id := 0

	if row.Next() {
		if err = row.Scan(&id); err != nil {
			log.Println(err)
			code, err = http.StatusInternalServerError, errors.New("internal server error")
			return
		}
	}

	if id != 0 {
		code, err = http.StatusConflict, errors.New("phone already used")
		log.Println(err)
		return
	}

	return
}
