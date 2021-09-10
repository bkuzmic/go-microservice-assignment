package app

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
			BadRequest(w, "Invalid request")
			return
		}
		var person models.Person
		err = json.Unmarshal(body, &person)
		if err != nil {
			log.Println("Error unmarshalling body request")
			log.Println(err)
			BadRequest(w, "Invalid request")
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

func (a *App) GetPersonHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			log.Println("ID parameter is missing")
			BadRequest(w, "ID parameter is missing")
		}
		person, err := a.DB.GetPerson(r.Context(), id)
		if err != nil {
			if err.Error() == "redis: nil" {
				NotFoundResponse(w)
			} else {
				ServerError(w)
			}
			return
		}
		OkResponse(w, person)
	}
}

func (a *App) UpdatePersonOptimisticHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// validate input
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error processing body request")
			log.Println(err)
			BadRequest(w, "Invalid request")
			return
		}
		var person models.Person
		err = json.Unmarshal(body, &person)
		if err != nil {
			log.Println("Error unmarshalling body request")
			log.Println(err)
			BadRequest(w, "Invalid request")
			return
		}
		if person.Id == "" {
			msg := "Missing person ID"
			log.Println(msg)
			BadRequest(w, msg)
			return
		}

		var modifiedPerson models.Person
		modifiedPerson, err = a.DB.UpdatePersonOptimistic(r.Context(), &person)
		OkResponse(w, &modifiedPerson)
	}
}

func (a *App) UpdatePersonPessimisticHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//person, err := a.DB.GetPerson(r.Context(), id)
		//if err != nil {
		//	if err.Error() == "redis: nil" {
		//		NotFoundResponse(w)
		//	} else {
		//		ServerError(w)
		//	}
		//	return
		//}
		//OkResponse(w, person)
	}
}

func ServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func NotFoundResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not found"))
}

func CreatedResponse(w http.ResponseWriter, person *models.Person) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	res, _ := json.Marshal(person)
	w.Write(res)
}

func OkResponse(w http.ResponseWriter, person *models.Person) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	res, _ := json.Marshal(person)
	w.Write(res)
}