package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"jeisaRaja.git/snippetbox/pkg/models"
)

func (app *application) signUpForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is the sign up form")
	return
}
func (app *application) signUp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is the sign up post request")
	return
}
func (app *application) logInForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is the log in form")
	return
}
func (app *application) logIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "this is the log in post")
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
	app.render(w, r, "home.page.tmpl", data)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) create_snippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
	}
	inputErr := make(map[string]string)
	w.Header().Set("Allow", "POST") // This is for adding another key-value pair to the header
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")
	if strings.TrimSpace(title) == "" {
		inputErr["Title"] = "Title cannot be empty!"
	} else if utf8.RuneCountInString(title) > 100 {
		inputErr["Title"] = "Maximum length for title is 100 characters"
	}
	if strings.TrimSpace(content) == "" {
		inputErr["Content"] = "Content cannot be empty!"
	}
	if strings.TrimSpace(expires) == "" {
		inputErr["Expires"] = "Expire date cannot be empty!"
	} else if expires != "365" && expires != "7" && expires != "1" {
		inputErr["Expires"] = "Expire date need to be one 1 ,7 or 365 days"
	}

	if len(inputErr) > 0 {
		data := &templateData{
			FormError: inputErr,
			FormData:  r.PostForm,
		}
		fmt.Println(data)
		app.render(w, r, "create.page.tmpl", data)
		return
	}
	id, err := app.snippets.Insert(title, content, expires)
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

	app.render(w, r, "create.page.tmpl", &templateData{})

}
