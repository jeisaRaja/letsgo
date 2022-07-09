package main

import (
	"fmt"
	"net/http"
)

func serverError() {
	fmt.Print()
}

func clienError() {
	fmt.Fprint()
	http.HandleFunc()
}
