package main

import (
	"net/http"
	"strconv"

	"jeisaRaja.git/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// fmt.Print(r.Method)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
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
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST") // This is for adding another key-value pair to the header
		title := "Fortune Cookie"
		content := "Fortune cookie, berbentuk hati, hey hey hey"
		expires := "7"
		app.snippets.Insert(title, content, expires)
		app.clienError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write(([]byte("Create a snippet...")))
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
