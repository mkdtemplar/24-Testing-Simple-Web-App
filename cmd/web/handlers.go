package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

type TemplateData struct {
	IP   string
	Data map[string]any
}

var pathToTemplates = "./templates/"

func (a *application) Home(w http.ResponseWriter, r *http.Request) {
	var td = make(map[string]any)

	if a.Session.Exists(r.Context(), "test") {
		msg := a.Session.GetString(r.Context(), "test")
		td["test"] = msg
	} else {
		a.Session.Put(r.Context(), "test", "Hit this page at "+time.Now().UTC().String())
	}
	_ = a.render(w, r, "home.page.gohtml", &TemplateData{Data: td})
}

func (a *application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// validate form
	form := NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		a.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := a.DB.GetUserByEmail(email)
	if err != nil {
		log.Println(err)
	}

	log.Println("From database: ", user.FirstName)
	log.Println(email, password)
	_, err = fmt.Fprintf(w, email)
	if err != nil {
		return
	}
}

func (a *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	// parse template from disk
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	data.IP = a.ipFromContext(r.Context())

	//execute template, passing data if any
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
