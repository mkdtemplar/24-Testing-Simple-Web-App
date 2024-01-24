package main

import (
	"24-Testing-Simple-Web-App/pkg/db"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	Session *scs.SessionManager
	DB      db.PostgresConn
	DSN     string
}

func main() {
	// Set up application configuration
	app := application{}

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.Parse()

	conn, err := app.connectToDb()
	if err != nil {
		log.Fatal(err)
	}

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	app.DB = db.PostgresConn{DB: conn}

	// get session manager
	app.Session = getSession()

	// Get application routes
	mux := app.routes()

	// Prion out the message
	fmt.Println("Starting server on port 8080....")

	// Start the server
	err = http.ListenAndServe("localhost:8080", mux)

	if err != nil {
		log.Fatalf("Failed to start server %s", err)
	}
}
