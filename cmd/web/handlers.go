package main

import (
	"fmt"
	"net/http"
)

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
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	email := r.PostForm.Get("email")

	_, err = app.user.Insert(username, password, email)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("User created ", "username", username)
	http.Redirect(w, r, fmt.Sprint("/"), http.StatusSeeOther)
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

	app.logger.Info("User will log in when it will be handled", "username", username)
	http.Redirect(w, r, fmt.Sprint("/"), http.StatusSeeOther)
}
