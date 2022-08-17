package main

import (
	"fmt"
	"net/http"
	"strconv"

	"jeisaRaja.git/snippetbox/pkg/forms"
	"jeisaRaja.git/snippetbox/pkg/models"
)

func (app *application) signUpForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
	return
}
func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.MatchesPattern("email")
	form.MinLength("password", 10)
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	}
	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		app.serverError(w, err)
	}

	app.session.Put(r, "flash", "Your signup was successfull")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}
func (app *application) logInForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
	return
}
func (app *application) logIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Println(err)
	}
	form := forms.New(
		r.PostForm,
	)
	form.Required("email", "password")
	if len(form.Errors) > 0 {
		fmt.Println(form.Errors, "error brok")
		return
	}
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	fmt.Println("idnya bro:", id)
	if err == models.ErrInvalidCredentials {
		form.Errors.Add("generic", "Email or Password is incorrect")
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "userID", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)

	return
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// fmt.Print(r.Method)
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: s}
	// fmt.Fprint(w, data)
	app.render(w, r, "home.page.tmpl", data)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	return
}

func (app *application) create_snippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
	}
	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Snipet Successfully Created...")
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil {
		app.notFound(w, 404)
		return
	}
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w, 404)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{
		Snippet: s,
	}
	app.render(w, r, "show.page.tmpl", data)
}

func (app *application) create_snippet_form(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}
