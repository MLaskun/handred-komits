package models

import (
	"database/sql"
	"errors"
)

type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username string, password string, email string) (int, error) {
	stmt := `INSERT INTO users (username, password, email)
	VALUES(?,?,?)`

	result, err := m.DB.Exec(stmt, username, password, email)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Authenticate(username string, password string) (User, error) {
	stmt := `SELECT username, password FROM users
	WHERE username = ? AND password = ?`

	row := m.DB.QueryRow(stmt, username, password)

	var u User

	err := row.Scan(&u.Username, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errors.New("No matching user currently registered")
		} else {
			return User{}, err
		}
	}
	return u, nil
}
