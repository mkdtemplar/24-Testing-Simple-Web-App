package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func openDb(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (a *application) connectToDb() (*sql.DB, error) {
	connection, err := openDb(a.DSN)

	if err != nil {
		return nil, err
	}

	log.Println("Connected to postgres")

	return connection, nil
}
