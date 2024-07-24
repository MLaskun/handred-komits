package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/MLaskun/handred-komits/internal/validator"
)

type userSignupForm struct {
	Username            string `form:"username"`
	Password            string `form:"password"`
	Email               string `form:"email"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	data := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {

}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	_, err = app.user.Authenticate(username, password)
	if err != nil {
		app.logger.Info("No such user registered yet")
		http.Redirect(w, r, fmt.Sprint("/"), http.StatusSeeOther)
		return
	}
	app.logger.Info("Found user", "username", username, "password", password)

	app.logger.Info("User will log in when it will be handled", "username", username)
	http.Redirect(w, r, fmt.Sprint("/"), http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {

}
func (app *application) purgeDatabase(w http.ResponseWriter, r *http.Request) {
	err := purgeDatabase(app.db)
	if err != nil {
		http.Error(w, "Failed to purge database", http.StatusInternalServerError)
		app.logger.Error(err.Error())
		return
	}

	err = createTables(app.db)
	if err != nil {
		http.Error(w, "Failed to create tables", http.StatusInternalServerError)
		app.logger.Error(err.Error())
		return
	}

	w.Write([]byte("Database reset successfully"))
	app.logger.Info("Database reset successfully")
	return
}

func purgeDatabase(db *sql.DB) error {
	_, err := db.Exec(`DROP TABLE IF EXISTS users, sessions;`)
	return err
}
func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE users (
			id INTEGER AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			email VARCHAR(255) NOT NULL,
			hashed_password CHAR(60) NOT NULL,
			created DATETIME NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE TABLE sessions (
			token CHAR(43) PRIMARY KEY,
			data BLOB NOT NULL,
			expiry TIMESTAMP(6) NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		CREATE INDEX session_expiry_idx ON sessions (expiry);
	`)
	return err
}
