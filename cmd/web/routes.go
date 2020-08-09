package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	// Middleware chain to add headers, manage panic error and logging every call
	myChain := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicChain := alice.New(app.session.Enable, noSurf, app.authenticate)

	router := mux.NewRouter()
	router.Handle("/", dynamicChain.ThenFunc(app.handHome)).Methods(http.MethodGet, http.MethodHead)
	router.Handle("/snippet/create", dynamicChain.Append(app.requireAuth).ThenFunc(app.createSnippetForm)).Methods(http.MethodGet, http.MethodHead)
	router.Handle("/snippet/create", dynamicChain.Append(app.requireAuth).ThenFunc(app.createSnippet)).Methods(http.MethodPost)
	router.Handle("/snippet/{id:[0-9]+}", dynamicChain.ThenFunc(app.showSnippet)).Methods(http.MethodGet, http.MethodHead)
	router.Handle("/user/signup", dynamicChain.ThenFunc(app.signupUserForm)).Methods(http.MethodGet)
	router.Handle("/user/signup", dynamicChain.ThenFunc(app.signupUser)).Methods(http.MethodPost)
	router.Handle("/user/login", dynamicChain.ThenFunc(app.loginUserForm)).Methods(http.MethodGet)
	router.Handle("/user/login", dynamicChain.ThenFunc(app.loginUser)).Methods(http.MethodPost)
	router.Handle("/user/logout", dynamicChain.Append(app.requireAuth).ThenFunc(app.logoutUser)).Methods(http.MethodPost)

	// Dir with static assets
	fileServ := http.FileServer(http.Dir("./ui/static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", fileServ))

	return myChain.Then(router)
}
