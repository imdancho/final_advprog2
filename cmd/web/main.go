// package main

// import (
// 	"context"
// 	"crypto/tls"
// 	"database/sql"
// 	"flag"
// 	"fmt"
// 	"html/template"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"aitunews.kz/snippetbox/pkg/models/mysql"

// 	_ "github.com/go-sql-driver/mysql"
// 	"github.com/golangcollege/sessions" // New import	fkmm flqo zpha zddd
// 	_ "github.com/lib/pq"
// 	"golang.org/x/time/rate"
// )

// type application struct {
// 	errorLog      *log.Logger
// 	infoLog       *log.Logger
// 	session       *sessions.Session
// 	snippets      *mysql.SnippetModel
// 	services      *mysql.ServiceModel
// 	appointments  *mysql.AppointmentModel
// 	templateCache map[string]*template.Template
// 	users         *mysql.UserModel
// }

// func main() {

// 	//Rate Limiting ------------------------------------------------
// 	// Create a limiter that allows 1 event per second.
// 	limiter := rate.NewLimiter(rate.Limit(1), 1)

// 	// Create a context for the limiter.
// 	ctx := context.Background()

// 	ticker := time.NewTicker(time.Second)
// 	defer ticker.Stop()

// 	go func() {
// 		for {
// 			// Wait for permission from the limiter.
// 			if err := limiter.Wait(ctx); err != nil {
// 				fmt.Println("Rate limit exceeded")
// 			}
// 			// Perform some task
// 			//fmt.Println("Task executed at:", time.Now().Format(time.StampMilli))
// 		}
// 	}()

// 	// Output "hello"
// 	fmt.Println("hello")

// 	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Fatal("Failed to open log file:", err)
// 	}
// 	defer file.Close()

// 	// Create a multiwriter that writes to both the file and os.Stdout (terminal)
// 	multiWriter := io.MultiWriter(file, os.Stdout)

// 	log.SetOutput(multiWriter)

// 	log.Println("\n\n---------------------------------\n")

// 	// dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
// 	// dsn := flag.String("dsn", "postgres://postgres:qwerty@localhost/snippetbox?sslmode=disable", "PostgreSQL data source name")
// 	dsn := flag.String("dsn", "postgres://imdancho:Ac4YWIZbWO1u8yHYE9bwA8q8xRUVbsbe@dpg-colunr20si5c73faeff0-a.singapore-postgres.render.com:5432/snippetbox_rpcq?sslmode=require", "PostgreSQL data source name")

// 	addr := flag.String("addr", ":4000", "HTTP network address")
// 	// addr := flag.String("addr", "", "HTTP network address")
// 	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
// 	flag.Parse()
// 	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
// 	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
// 	infoLog.SetOutput(multiWriter)
// 	errorLog.SetOutput(multiWriter)

// 	db, err := openDB(*dsn)
// 	if err != nil {
// 		errorLog.Fatal(err)
// 	}
// 	defer db.Close()
// 	templateCache, err := newTemplateCache("./ui/html/")
// 	if err != nil {
// 		errorLog.Fatal(err)
// 	}

// 	session := sessions.New([]byte(*secret))
// 	session.Lifetime = 12 * time.Hour
// 	session.Secure = true
// 	// Initialize a mysql.UserModel instance and add it to the application
// 	// dependencies.
// 	app := &application{
// 		errorLog:      errorLog,
// 		infoLog:       infoLog,
// 		session:       session,
// 		snippets:      &mysql.SnippetModel{DB: db},
// 		services:      &mysql.ServiceModel{DB: db},
// 		appointments:  &mysql.AppointmentModel{DB: db},
// 		templateCache: templateCache,
// 		users:         &mysql.UserModel{DB: db},
// 	}
// 	tlsConfig := &tls.Config{
// 		PreferServerCipherSuites: true,
// 		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
// 	}
// 	srv := &http.Server{
// 		Addr:         *addr,
// 		ErrorLog:     errorLog,
// 		Handler:      app.routes(),
// 		TLSConfig:    tlsConfig,
// 		IdleTimeout:  time.Minute,
// 		ReadTimeout:  5 * time.Second,
// 		WriteTimeout: 10 * time.Second,
// 	}
// 	infoLog.Printf("Starting server on %s", *addr)
// 	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
// 	errorLog.Fatal(err)

// 	//--------------------------------------
// 	select {}
// }

// // The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// // The openDB() function wraps sql.Open() and returns a sql.DB connection pool
// // for a given DSN.

// // func openDB(dsn string) (*sql.DB, error) {
// // 	db, err := sql.Open("mysql", dsn)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	if err = db.Ping(); err != nil {
// // 		return nil, err
// // 	}
// // 	return db, nil
// // }

// func openDB(dsn string) (*sql.DB, error) {
// 	// Attempt to open a database connection using the provided DSN (Data Source Name)
// 	db, err := sql.Open("postgres", dsn)
// 	if err != nil {
// 		// If there was an error opening the connection, return nil for the database connection
// 		// and the error encountered.
// 		return nil, err
// 	}

// 	// Attempt to ping the database to ensure the connection is valid and the database is accessible.
// 	if err = db.Ping(); err != nil {
// 		// If there was an error pinging the database, close the database connection and return
// 		// nil for the database connection and the error encountered.
// 		db.Close() // Close the database connection before returning.
// 		return nil, err
// 	}

// 	// If everything was successful, return the database connection and nil for the error.
// 	return db, nil
// }

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"aitunews.kz/snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	_ "github.com/lib/pq"
	"golang.org/x/time/rate"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	snippets      *mysql.SnippetModel
	services      *mysql.ServiceModel
	appointments  *mysql.AppointmentModel
	templateCache map[string]*template.Template
	users         *mysql.UserModel
}

func main() {
	//Rate Limiting
	limiter := rate.NewLimiter(rate.Limit(1), 1)
	ctx := context.Background()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		for {
			if err := limiter.Wait(ctx); err != nil {
				fmt.Println("Rate limit exceeded")
			}
		}
	}()

	fmt.Println("hello")

	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer file.Close()

	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multiWriter)

	log.Println("\n\n---------------------------------\n")

	dsn := flag.String("dsn", "postgres://imdancho:Ac4YWIZbWO1u8yHYE9bwA8q8xRUVbsbe@dpg-colunr20si5c73faeff0-a.singapore-postgres.render.com:5432/snippetbox_rpcq?sslmode=require", "PostgreSQL data source name")
	addr := flag.String("addr", ":4000", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog.SetOutput(multiWriter)
	errorLog.SetOutput(multiWriter)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = false // Set to false for HTTP
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		snippets:      &mysql.SnippetModel{DB: db},
		services:      &mysql.ServiceModel{DB: db},
		appointments:  &mysql.AppointmentModel{DB: db},
		templateCache: templateCache,
		users:         &mysql.UserModel{DB: db},
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
