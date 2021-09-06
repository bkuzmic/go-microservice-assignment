package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const LayoutDateOnly = "2006-02-01"

type JSONDate time.Time

func (date JSONDate) MarshalJSON() ([]byte, error) {
	d := fmt.Sprintf("\"%s\"", time.Time(date).Format(LayoutDateOnly))
	return []byte(d), nil
}

type Person struct {
	Name string `json:"Name"`
	Address string `json:"Address"`
	DateOfBirth JSONDate `json:"DateOfBirth"`
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint Hit: GetPerson")
	var person = Person{
		Name:        "Boris",
		Address:     "Zagreb",
		DateOfBirth: JSONDate(time.Date(1981, time.November, 29, 0, 0, 0, 0, time.UTC)),
	}
	json.NewEncoder(w).Encode(person)
}

func handleRequests() {
	http.HandleFunc("/person", GetPerson)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func main() {
	handleRequests()
}
