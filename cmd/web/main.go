package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mhmtnzly/bookings/pkg/config"
	"github.com/mhmtnzly/bookings/pkg/handlers"
	"github.com/mhmtnzly/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	// change this to true when in production
	app.Inproduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.Inproduction

	app.Session = session
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	render.NewTemplates(&app)

	// http.HandleFunc("/", handlers.Repo.Home)
	// http.HandleFunc("/about", handlers.Repo.About)
	fmt.Printf("Staring application on port %s", portNumber)
	// _ = http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
