package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/thats-insane/awt-quiz3/internal/validator"
)

/* Common error
 */
var ErrRecordNotFound = errors.New("record not found")

/* User struct
 */
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

/* User database model struct
 */
type UserModel struct {
	DB *sql.DB
}

/* Insert a user into the database
 */
func (u UserModel) Insert(user *User) error {
	qry := `
	INSERT INTO users (fullname, email) 
	VALUES ($1, $2) 
	RETURNING id, created_at, version`

	args := []any{user.Name, user.Email}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, qry, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
}

/* Locate and return a user from the database
 */
func (u UserModel) Get(id int64) (*User, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	qry := `
	SELECT id, fullname, email, created_at, version 
	FROM user 
	WHERE id = $1`

	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, qry, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

/* Update a user in the database
 */
func (u UserModel) Update(user *User) error {
	qry := `
	UPDATE user 
	SET fullname = $1, email = $2, version = version + 1 
	WHERE id = $3 
	RETURNING version`

	args := []any{user.Name, user.Email, user.ID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return u.DB.QueryRowContext(ctx, qry, args...).Scan(&user.Version)
}

/* Delete a user from the database
 */
func (u UserModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	qry := `
	DELETE FROM user 
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := u.DB.ExecContext(ctx, qry, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func Validate(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "fullname", "cannot be empty")
	v.Check(user.Email != "", "email", "cannot be empty")
	v.Check(len(user.Name) <= 25, "fullname", "cannot be more than 25 bytes long")
	v.Check(len(user.Email) <= 50, "email", "cannot be more than 50 bytes long")
}
