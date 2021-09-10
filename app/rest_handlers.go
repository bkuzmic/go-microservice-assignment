package app

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go-microservice-assignment/app/models"
	"io/ioutil"
	"log"
	"net/http"
)

func (a *App) IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Microservices Assignment v.0.0.1")
	}
}

func (a *App) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (a *App) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (a *App) CreatePersonHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error processing body request")
			log.Println(err)
			BadRequest(w)
			return
		}
		var person models.Person
		err = json.Unmarshal(body, &person)
		if err != nil {
			log.Println("Error unmarshalling body request")
			log.Println(err)
			BadRequest(w)
			return
		}

		key := uuid.New().String()
		person.Id = key

		err = a.DB.CreatePerson(r.Context(), &person)
		if err != nil {
			log.Println("Error writing person to storage")
			log.Println(err)
			ServerError(w)
			return
		}

		CreatedResponse(w, &person)
	}
}

func ServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
}

func CreatedResponse(w http.ResponseWriter, person *models.Person) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res, _ := json.Marshal(person)
	w.Write(res)
}