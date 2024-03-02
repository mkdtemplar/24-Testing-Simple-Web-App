package main

import (
	"24-Testing-Simple-Web-App/pkg/data"
	"24-Testing-Simple-Web-App/pkg/repository"
	"24-Testing-Simple-Web-App/pkg/repository/dbrepo"
	"database/sql"
	"encoding/gob"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

const port = "8090"

type application struct {
	DSN       string
	DB        repository.DataBaseRepo
	Domain    string
	JWTSecret string
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	gob.Register(data.User{})
	// Set up application configuration
	app := application{}
	flag.StringVar(&app.Domain, "domain", "example.com", "Domain for application, e.g. company.com")
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160", "signing secret")

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

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	log.Printf("Starting APi on port %s\n", port)

	err = app.routes().Run(fmt.Sprintf("localhost:%s", port))
	if err != nil {
		fmt.Println(err)
		return
	}

}
