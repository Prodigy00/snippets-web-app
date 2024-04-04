package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/Prodigy00/snippetbox/internal/models"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

// Don't use DefaultServeMux, always define a http.NewServeMux() in main
// third party packages can write to DefaultServeMux. This would cause issues.
func main() {
	//cli flags preferred to env vars for addr config
	addr := flag.String("addr", ":4001", "HTTP network address")
	sqlUser := os.Getenv("SQL_USER_NAME")
	sqlPwd := os.Getenv("SQL_PWD")

	dsnVal := fmt.Sprintf("%s:%s@/snippetbox?parseTime=true", sqlUser, sqlPwd)

	dsn := flag.String("dsn", dsnVal, "MySQL data source name")
	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are // encountered during parsing the application will be terminated.
	flag.Parse()

	// Use log.New() to create a logger for writing information messages. This takes
	// three parameters: the destination to write the logs to (os.Stdout), a string // prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flags // are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as
	//the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	//that the server uses the same network address and routes as before, and set
	//the ErrorLog field so that the server now uses the custom errorLog logger in
	//the event of any problems.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}
	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe() // eq to app.listen()
	errorLog.Fatal(err)
}

// cmd - app-specific code(business logic)
// internal - non-app-specific code e.g validation helpers and SQL db models
//ui

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
