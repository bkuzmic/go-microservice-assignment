package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-microservice-assignment/app"
	"go-microservice-assignment/app/storage"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var rdb *redis.Client
var ctx = context.Background()

func main() {
	log.Println("Connecting to Redis database...")
	check(connectToRedis())

	keyExpireTime, err := strconv.Atoi(os.Getenv("KEY_IDLE_TIME_MINUTES"))
	check(err)

	db := storage.NewDB(rdb, time.Duration(keyExpireTime)*time.Minute)

	application := app.New(db)
	http.HandleFunc("/", application.Router.ServeHTTP)

	log.Println("Application started at port 8000")
	err = http.ListenAndServe(":8000", nil)
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
