package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/nosurf"
)

func (app *application) authenticateUser(r *http.Request) int {
	return app.session.GetInt(r, "userID")
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(3, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clienError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	timeNow := time.Now().Year()
	if td == nil {
		td = &templateData{}
	}
	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticateUser(r)
	fmt.Println(td.AuthenticatedUser)
	td.CurrentYear = timeNow
	td.Flash = app.session.PopString(r, "flash")
	return td

}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	buffer := new(bytes.Buffer)
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("The template of %s error", name))
		return
	}
	fmt.Println(td)
	err := ts.Execute(buffer, app.addDefaultData(td, r))
	if err != nil {
		fmt.Println(err)
		return
	}
	buffer.WriteTo(w)
}
