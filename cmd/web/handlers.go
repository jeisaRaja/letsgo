package main

import (
	"fmt"
	"net/http"
	"strconv"

	"jeisaRaja.git/snippetbox/pkg/models"
)

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
	w.Header().Set("Allow", "POST") // This is for adding another key-value pair to the header
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
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
	}
	data := &templateData{Snippet: s}
	app.render(w, r, "show.page.tmpl", data)
}

func (app *application) create_snippet_form(w http.ResponseWriter, r *http.Request) {

	app.render(w, r, "create_form.page.tmpl", &templateData{})

}
