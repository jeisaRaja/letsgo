package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// fmt.Print(r.Method)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Display the home page"))
}
func snippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet"))
}
func create_snippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST") // This is for adding another key-value pair to the header
		w.WriteHeader(405)
		w.Write([]byte("This method is not supported \n"))
		return
	}
	w.Write(([]byte("Create a snippet...")))
}
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
