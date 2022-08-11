package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"jeisaRaja.git/snippetbox/pkg/models/mysql"
)

type application struct {
	infoLog       *log.Logger
	errorLog      *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	session       *sessions.Session
}

func main() {

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	defaultDSN := "web:mysqlCirebon01@/snippetbox?parseTime=true"
	addr := flag.String("addr", "4000", "HTTP network address")
	dsn := flag.String("dsn", defaultDSN, "MYSQL Database Pool")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key for session")
	flag.Parse()

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	db, err := openDB(*dsn)
	if err != nil {
		errlog.Println("DB Connection Error!")
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errlog.Fatal()
	}

	app := &application{
		infoLog:       infoLog,
		errorLog:      errlog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
	}
	*addr = "127.0.0.1:" + *addr
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errlog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
