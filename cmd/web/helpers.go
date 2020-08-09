package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/eze8789/snippetbox/pkg/models"
	"github.com/justinas/nosurf"
)

func (app *application) serveErr(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientErr(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientErr(w, http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	td.AuthenticatedUser = app.authUser(r)
	td.CSRFToken = nosurf.Token(r)
	return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serveErr(w, fmt.Errorf("The template %s does not exist", name))
		return
	}
	// Initialize buffer to write before sent data back
	b := new(bytes.Buffer)

	err := ts.Execute(b, app.addDefaultData(td, r))
	if err != nil {
		app.serveErr(w, err)
		return
	}
	b.WriteTo(w)
}

func (app *application) authUser(r *http.Request) *models.User {
	u, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}

	return u
}
