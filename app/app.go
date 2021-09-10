package app

import (
	"github.com/gorilla/mux"
	"go-microservice-assignment/app/storage"
)

type App struct {
	Router *mux.Router
	DB *storage.DB
}

func New(db *storage.DB) *App {
	app:= &App {
		Router: mux.NewRouter(),
		DB: db,
	}
	app.initRoutes()
	return app
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/", a.IndexHandler()).Methods("GET")
	a.Router.HandleFunc("/health", a.HealthHandler()).Methods("GET")
	a.Router.HandleFunc("/readiness", a.ReadinessHandler()).Methods("GET")
	a.Router.HandleFunc("/api/v1/person", a.CreatePersonHandler()).Methods("POST")
	a.Router.HandleFunc("/api/v1/person/{id}", a.GetPersonHandler()).Methods("GET")
	a.Router.HandleFunc("/api/v1/person", a.UpdatePersonOptimisticHandler()).Methods("PATCH")
	a.Router.HandleFunc("/api/v1/person/pessimistic", a.UpdatePersonPessimisticHandler()).Methods("PATCH")
}