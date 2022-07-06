package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")

	flag.Parse()

	log.Println(*addr)
	// The flag parse didnt work on ubuntu virtualbox, there is nothing with the code as in windows, it works perfectly
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", snippet)
	mux.HandleFunc("/snippet/create", create_snippet)
	fileserver := http.FileServer(http.Dir("../../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileserver))
	log.Println("Starting server on port", *addr)

	err := http.ListenAndServe(*addr, mux)
	log.Println(err)
}
