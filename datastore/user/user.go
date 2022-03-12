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
	// prepare statement before execution
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
	now := time.Now()
	res, err := stmt.Exec(newUser.Name, newUser.Email, newUser.Phone, newUser.Password, now, now)

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
	user.CreatedAt = now
	user.UpdatedAt = now

	return
}

func (ur *UserRepository) GetUserByEmail(email string) (user _entity.User, code int, err error) {
	// prepare statment before execution
	stmt, err := ur.db.Prepare(`
		SELECT id, role, name, phone, password, created_at, updated_at
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

	if row.Next() {
		if err = row.Scan(&user.Id, &user.Role, &user.Name, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Println(err)
			code, err = http.StatusInternalServerError, errors.New("internal server error")
			return
		}
	}

	return
}

func (ur *UserRepository) GetUserById(userId int) (user _entity.User, code int, err error) {
	// prepare statement before execution
	stmt, err := ur.db.Prepare(`
		SELECT role, name, email, phone, password, created_at, updated_at
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
		if err = row.Scan(&user.Role, &user.Name, &user.Email, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Println(err)
			code, err = http.StatusInternalServerError, errors.New("internal server error")
			return
		}
	}

	return
}

func (ur *UserRepository) UpdateUser(updatedUser _entity.User) (user _entity.User, code int, err error) {
	// prepare statement before execution
	stmt, err := ur.db.Prepare(`
		UPDATE users
		SET name = ?, email = ?, phone = ?, password = ?, updated_at = ?
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
	now := time.Now()
	_, err = stmt.Exec(updatedUser.Name, updatedUser.Email, updatedUser.Phone, updatedUser.Password, now, updatedUser.Id)

	if err != nil {
		log.Println(err)
		code, err = http.StatusInternalServerError, errors.New("internal server error")
		return
	}

	user = updatedUser
	user.UpdatedAt = now

	return
}

func (ur *UserRepository) DeleteUser(userId int) (code int, err error) {
	// prepare statement before execution
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
