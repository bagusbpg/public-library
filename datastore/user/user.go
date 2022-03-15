package user

import (
	"log"
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

func (ur *UserRepository) CreateNewUser(newUser _entity.User) (user _entity.User, err error) {
	// prepare statement
	stmt, err := ur.db.Prepare(`
		INSERT INTO users (role, name, email, phone, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	res, err := stmt.Exec(newUser.Role, newUser.Name, newUser.Email, newUser.Phone, newUser.Password, newUser.CreatedAt, newUser.UpdatedAt)

	if err != nil {
		log.Println(err)
		return
	}

	// get new user id
	id, err := res.LastInsertId()

	if err != nil {
		log.Println(err)
		return
	}

	user = newUser
	user.Id = int(id)

	return
}

func (ur *UserRepository) GetUserByEmail(email string) (user _entity.User, err error) {
	// prepare statment before execution
	stmt, err := ur.db.Prepare(`
		SELECT id, role, name, phone, password, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		  AND email = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	row, err := stmt.Query(email)

	if err != nil {
		log.Println(err)
		return
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&user.Id, &user.Role, &user.Name, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (ur *UserRepository) GetUserById(userId int) (user _entity.User, err error) {
	// prepare statement
	stmt, err := ur.db.Prepare(`
		SELECT role, name, email, phone, password, created_at, updated_at
		FROM users
		WHERE deleted_at IS NULL
		  AND id = ?
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

	if row.Next() {
		if err = row.Scan(&user.Role, &user.Name, &user.Email, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Println(err)
			return
		}
	}

	return
}

func (ur *UserRepository) UpdateUser(updatedUser _entity.User) (user _entity.User, err error) {
	// prepare statement
	stmt, err := ur.db.Prepare(`
		UPDATE users
		SET name = ?, email = ?, phone = ?, password = ?, updated_at = ?
		WHERE id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(updatedUser.Name, updatedUser.Email, updatedUser.Phone, updatedUser.Password, updatedUser.UpdatedAt, updatedUser.Id)

	if err != nil {
		log.Println(err)
		return
	}

	user = updatedUser

	return
}

func (ur *UserRepository) DeleteUser(userId int) (err error) {
	// prepare statement
	stmt, err := ur.db.Prepare(`
		UPDATE users
		SET deleted_at = ?
		WHERE id = ?
	`)

	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	// execute statement
	_, err = stmt.Exec(time.Now(), userId)

	if err != nil {
		log.Println(err)
		return
	}

	return
}
