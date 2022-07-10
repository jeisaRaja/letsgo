package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	addr := flag.String("addr", "4000", "HTTP network address")

	flag.Parse()

	// The flag parse didnt work on ubuntu virtualbox, there is nothing with the code as in windows, it works perfectly
	app := &application{
		infoLog:  infoLog,
		errorLog: errlog,
	}

	*addr = "127.0.0.1:" + *addr

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errlog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errlog.Fatal(err)
}
