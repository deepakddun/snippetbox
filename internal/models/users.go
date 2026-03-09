package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (m *UserModel) Insert(ctx context.Context, name, email, password string) error {

	var pgErr *pgconn.PgError
	stmt := `
	INSERT INTO users (name, email, hashed_password, created)
	VALUES ($1, $2, $3,NOW())
	RETURNING id
`

	hasshedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	var id int

	err = m.DB.QueryRow(ctx, stmt, name, email, hasshedPassword).Scan(&id)

	if errors.As(err, &pgErr) {
		fmt.Println(err)
		if pgErr.Code == "23505" {
			return ErrDuplicateEmail
		}

	}

	if err != nil {
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
