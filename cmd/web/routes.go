package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *application) routes() http.Handler {
	mux := chi.NewRouter()

	// Register middleware
	mux.Use(middleware.Recoverer)
	mux.Use(a.addIpToContext)
	mux.Use(a.Session.LoadAndSave)

	// Register routes
	mux.Get("/", a.Home)
	mux.Post("/login", a.Login)

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(a.auth)
		mux.Get("/profile", a.Profile)
		mux.Post("/upload-profile-pic", a.UploadProfilePicture)
	})

	// Static assets
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
