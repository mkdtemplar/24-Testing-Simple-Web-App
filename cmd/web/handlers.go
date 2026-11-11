package main

import (
	"24-Testing-Simple-Web-App/pkg/data"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

type TemplateData struct {
	IP    string
	Data  map[string]any
	Error string
	Flash string
	User  data.User
}

type UploadedFile struct {
	OriginalFileName string
	FileSize         int64
}

var pathToTemplates = "./templates/"
var uploadPath = "./static/img"

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

	if a.Session.Exists(r.Context(), "user") {
		td.User = a.Session.Get(r.Context(), "user").(data.User)
	}

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

func (a *application) UploadFiles(r *http.Request, uploadDir string) ([]*UploadedFile, error) {
	var uploadedFiles []*UploadedFile

	err := r.ParseMultipartForm(int64(5242880))
	if err != nil {
		return nil, fmt.Errorf("uploaded file is bigger than %d bytes", 5242880)
	}

	for _, fileHeaders := range r.MultipartForm.File {
		for _, header := range fileHeaders {
			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {
				uploadedFile := UploadedFile{}
				var outfile *os.File
				inFile, err := header.Open()
				if err != nil {
					return nil, err
				}
				defer func(inFile multipart.File) {
					err := inFile.Close()
					if err != nil {
						return
					}
				}(inFile)
				uploadedFile.OriginalFileName = header.Filename
				defer func(outfile *os.File) {
					err := outfile.Close()
					if err != nil {
						return
					}
				}(outfile)
				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.OriginalFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, inFile)
					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}
				uploadedFiles = append(uploadedFiles, &uploadedFile)
				return uploadedFiles, nil
			}(uploadedFiles)
			if err != nil {
				return uploadedFiles, err
			}
		}
	}
	return uploadedFiles, nil
}

func (a *application) UploadProfilePicture(writer http.ResponseWriter, request *http.Request) {

	files, err := a.UploadFiles(request, uploadPath)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	user := a.Session.Get(request.Context(), "user").(data.User)

	img := data.UserImage{
		UserID:   user.ID,
		FileName: files[0].OriginalFileName,
	}

	_, err = a.DB.InsertUserImage(img)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	updatedUser, err := a.DB.GetUser(user.ID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	a.Session.Put(request.Context(), "user", updatedUser)

	http.Redirect(writer, request, "/user/profile", http.StatusSeeOther)
}
