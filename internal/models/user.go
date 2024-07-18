package models

import "database/sql"

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
