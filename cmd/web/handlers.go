package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/eze8789/snippetbox/pkg/forms"
	"github.com/eze8789/snippetbox/pkg/models"
	"github.com/gorilla/mux"
)

func (app *application) handHome(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serveErr(w, err)
		return
	}

	app.renderTemplate(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serveErr(w, err)
	}

	app.renderTemplate(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientErr(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	// Manage possible errors of the form data
	form.Required("title", "content", "expires")
	form.Length("title", 5, 30)
	form.PermittedValues("expires", "365", "7", "1")
	if !form.Valid() {
		app.renderTemplate(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	// Insert content in the db based on the received data
	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serveErr(w, err)
		return
	}

	app.session.Put(r, "flash", "Snippet created!")

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientErr(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	// Manage possible errors of the form data
	form.Required("name", "email", "password")
	form.MatchPattern("email", forms.EmailRegExp)
	form.MinLength("password", 10)
	if !form.Valid() {
		app.renderTemplate(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	err = app.user.InsertUser(form.Get("name"), form.Get("email"), form.Get("password"))
	if err == models.ErrDupEmail {
		form.Errors.Add("email", "Email address already in use")
		app.renderTemplate(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serveErr(w, err)
		return
	}

	app.session.Put(r, "flash", "User created")
	http.Redirect(w, r, fmt.Sprintf("/user/login"), http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.renderTemplate(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientErr(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		app.renderTemplate(w, r, "login.page.tmpl", &templateData{Form: form})
	}

	id, err := app.user.Authenticate(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCreds {
		form.Errors.Add("generic", "Invalid credentials")
		app.renderTemplate(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serveErr(w, err)
	}
	app.session.Put(r, "userID", id)

	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "userID")
	app.session.Put(r, "flash", "Logged Out")

	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}
