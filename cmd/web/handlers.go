package main

import (
	"fmt"
	"html/template"
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
	// array of templates:
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		app.errorLog.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func (app *application) snippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.infoLog.Println("No ID found")
		app.notFound(w, 404)
		return
	}
	s, err := app.snippets.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w, 500)
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Fprintf(w, "%v", s)
	// w.Write([]byte("Display a specific snippet"))
}

func (app *application) create_snippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST") // This is for adding another key-value pair to the header

		// w.WriteHeader(405)
		// w.Write([]byte("This method is not supported \n"))
		title := "Fortune Cookie"
		content := "Fortune cookie, berbentuk hati, hey hey hey"
		expires := "7"
		app.snippets.Insert(title, content, expires)
		app.clienError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write(([]byte("Create a snippet...")))
}
