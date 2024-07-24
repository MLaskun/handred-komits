package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int
	Username       string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username string, password string, email string) error {
	return nil
}

func (m *UserModel) Authenticate(username string, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
