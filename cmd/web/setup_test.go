package main

import (
	"24-Testing-Simple-Web-App/pkg/repository/dbrepo"
	"os"
	"testing"
)

var app application

func TestMain(m *testing.M) {
	pathToTemplates = "./../../templates/"
	app.Session = getSession()

	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run())

}
