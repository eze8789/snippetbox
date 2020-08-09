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

	"github.com/eze8789/snippetbox/pkg/models/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
)

// Application struct, so we can access dependencies through methods
type application struct {
	errLog        *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
	session       *sessions.Session
	user          *mysql.UserModel
}

type contextKey string

var contextKeyUser = contextKey("user")

func main() {
	// Retrieve address and port from cli arg
	addr := flag.String("addr", ":3000", "HTTP Address and port")

	// Mysql connection string
	connStr := flag.String("dsn", "USERNAME:PASSWORD@/snippetbox?parseTime=true", "MySQL data source name")

	// Flag to mannage session secret key
	secret := flag.String("secret", "SECRET_KEY", "Secret Key")
	flag.Parse()

	// Manage informational and error logs
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*connStr)
	if err != nil {
		errLog.Fatal(err)
	}

	// When main function exits db pool connectios are closed
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html")
	if err != nil {
		errLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// Initialize app with their dependencies
	app := &application{
		errLog:        errLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
		session:       session,
		user:          &mysql.UserModel{DB: db},
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		MinVersion: tls.VersionTLS12,
	}
	// Initialize webserver struct and user errLog for logging events.
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on port %s", *addr)
	// Call the server as a method of the struct
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	// err = srv.ListenAndServe()
	errLog.Fatal(err)
}

func openDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
