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
	log.Println("Starting server on port", port)
	err := http.ListenAndServe(port, mux)
	log.Println(err)
}
