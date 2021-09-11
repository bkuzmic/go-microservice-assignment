package app

import (
	"github.com/gorilla/mux"
	"go-microservice-assignment/app/storage"
)

type app struct {
	Router *mux.Router
	DB storage.RedisDB
}

func New(db storage.RedisDB) *app {
	app:= &app {
		Router: mux.NewRouter(),
		DB: db,
	}
	app.initRoutes()
	return app
}

func (a *app) initRoutes() {
	a.Router.HandleFunc("/", a.IndexHandler()).Methods("GET")
	a.Router.HandleFunc("/health", a.HealthHandler()).Methods("GET")
	a.Router.HandleFunc("/readiness", a.ReadinessHandler()).Methods("GET")
	a.Router.HandleFunc("/api/v1/person", a.CreatePersonHandler()).Methods("POST")
	a.Router.HandleFunc("/api/v1/person/{id}", a.GetPersonHandler()).Methods("GET")
	a.Router.HandleFunc("/api/v1/person", a.UpdatePersonOptimisticHandler()).Methods("PATCH")
	a.Router.HandleFunc("/api/v1/person/pessimistic", a.UpdatePersonPessimisticHandler()).Methods("PATCH")
}