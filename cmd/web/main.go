package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", snippet)
	mux.HandleFunc("/snippet/create", create_snippet)
	port := ":4000"
	fileserver := http.FileServer(http.Dir("../../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	err := http.ListenAndServe(port, mux)
	log.Println("Starting server on port", port)
	log.Println(err)
}
