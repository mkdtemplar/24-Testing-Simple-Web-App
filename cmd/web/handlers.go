package main

import (
	"24-Testing-Simple-Web-App/pkg/data"
	"html/template"
	"log"
	"net/http"
	"path"
	"time"
)

type TemplateData struct {
	IP    string
	Data  map[string]any
	Error string
	Flash string
	User  data.User
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

func (a *application) Profile(w http.ResponseWriter, r *http.Request) {

	_ = a.render(w, r, "profile.page.gohtml", &TemplateData{})
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
		a.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// authenticate the user, if not redirect with error
	if !a.authenticate(r, user, password) {
		a.Session.Put(r.Context(), "error", "Invalid login")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// prevent fixation attack
	_ = a.Session.RenewToken(r.Context())

	// redirect to some other page
	a.Session.Put(r.Context(), "flash", "Successfully logged in!")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)
}

func (a *application) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {
	// parse template from disk
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	td.IP = a.ipFromContext(r.Context())
	td.Error = a.Session.PopString(r.Context(), "error")
	td.Flash = a.Session.PopString(r.Context(), "flash")

	//execute template, passing data if any
	err = parsedTemplate.Execute(w, td)
	if err != nil {
		return err
	}
	return nil
}

func (a *application) authenticate(r *http.Request, user *data.User, password string) bool {
	if valid, err := user.PasswordMatches(password); err != nil || !valid {
		return false
	}

	a.Session.Put(r.Context(), "user", user)
	return true
}
