package app

import (
	"encoding/json"
	"fmt"
	"go-microservice-assignment/app/models"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (a *app) IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go Person Service v.0.0.1")
	}
}

func (a *app) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (a *app) ReadinessHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (a *app) CreatePersonHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error processing body request:", err)
			badRequest(w, "Invalid request")
			return
		}
		var person models.Person
		err = json.Unmarshal(body, &person)
		if err != nil {
			log.Println("Error unmarshalling body request:", err)
			badRequest(w, "Invalid request")
			return
		}

		key := uuid.New().String()
		person.Id = key

		err = a.DB.CreatePerson(r.Context(), &person)
		if err != nil {
			log.Println("Error writing person to storage:", err)
			serverError(w)
			return
		}

		createdResponse(w, &person)
	}
}

func (a *app) GetPersonHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			log.Println("ID parameter is missing")
			badRequest(w, "ID parameter is missing")
			return
		}
		person, err := a.DB.GetPerson(r.Context(), id)
		if err != nil {
			if err.Error() == "redis: nil" {
				notFoundResponse(w)
			} else {
				serverError(w)
			}
			return
		}
		okResponse(w, person)
	}
}

func (a *app) UpdatePersonOptimisticHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// validate input
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error processing body request")
			log.Println(err)
			badRequest(w, "Invalid request")
			return
		}
		var person models.Person
		err = json.Unmarshal(body, &person)
		if err != nil {
			log.Println("Error unmarshalling body request")
			log.Println(err)
			badRequest(w, "Invalid request")
			return
		}
		if person.Id == "" {
			msg := "Missing person ID"
			log.Println(msg)
			badRequest(w, msg)
			return
		}

		modifiedPerson, err := a.DB.UpdatePersonOptimistic(r.Context(), &person)
		if err != nil {
			log.Println("Error while calling UpdatePersonOptimistic", err)
			serverError(w)
			return
		}
		okResponse(w, modifiedPerson)
	}
}

func (a *app) UpdatePersonPessimisticHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// validate input
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error processing body request")
			log.Println(err)
			badRequest(w, "Invalid request")
			return
		}
		var person models.Person
		err = json.Unmarshal(body, &person)
		if err != nil {
			log.Println("Error unmarshalling body request")
			log.Println(err)
			badRequest(w, "Invalid request")
			return
		}
		if person.Id == "" {
			msg := "Missing person ID"
			log.Println(msg)
			badRequest(w, msg)
			return
		}

		modifiedPerson, err := a.DB.UpdatePersonPessimistic(r.Context(), &person)
		if err != nil {
			log.Println("Error while calling UpdatePersonOptimistic", err)
			serverError(w)
			return
		}
		okResponse(w, modifiedPerson)
	}
}

func serverError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func badRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func notFoundResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not found"))
}

func createdResponse(w http.ResponseWriter, person *models.Person) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res, _ := json.Marshal(person)
	w.Write(res)
}

func okResponse(w http.ResponseWriter, person *models.Person) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(person)
	w.Write(res)
}
