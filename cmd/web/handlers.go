package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	// fmt.Print(r.Method)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ts, err := template.ParseFiles("../../ui/html/home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.Write([]byte("Display the home page"))
}
func snippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Displaying %d", id)
	// w.Write([]byte("Display a specific snippet"))
}
func create_snippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST") // This is for adding another key-value pair to the header
		// w.WriteHeader(405)
		// w.Write([]byte("This method is not supported \n"))
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write(([]byte("Create a snippet...")))
}
