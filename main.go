package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-microservice-assignment/app"
	"go-microservice-assignment/app/storage"
	"log"
	"net/http"
	"os"
)

//const LayoutDateOnly = "2006-02-01"
//
//type JSONDate time.Time
//
//func (date JSONDate) MarshalJSON() ([]byte, error) {
//	d := fmt.Sprintf("\"%s\"", time.Time(date).Format(LayoutDateOnly))
//	return []byte(d), nil
//}
//
//type Person struct {
//	Name string `json:"Name"`
//	Address string `json:"Address"`
//	DateOfBirth JSONDate `json:"DateOfBirth"`
//}
//
//func GetPerson(w http.ResponseWriter, r *http.Request) {
//	log.Println("Endpoint Hit: GetPerson")
//	var person = Person{
//		Name:        "Boris",
//		Address:     "Zagreb",
//		DateOfBirth: JSONDate(time.Date(1981, time.November, 29, 0, 0, 0, 0, time.UTC)),
//	}
//	json.NewEncoder(w).Encode(person)
//}
//

var rdb *redis.Client
var ctx = context.Background()

func main() {
	log.Println("Connecting to Redis database...")
	check(connectToRedis())

	storage := &storage.DB {
		Client: rdb,
	}

	app := app.New(storage)
	http.HandleFunc("/", app.Router.ServeHTTP)

	log.Println("Application started at port 8000")
	err := http.ListenAndServe(":8000", nil)
	check(err)
}

func connectToRedis() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	return err
}

func check(e error) {
	if e != nil {
		log.Println(e)
		os.Exit(1)
	}
}
