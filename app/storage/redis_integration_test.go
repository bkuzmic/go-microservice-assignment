// +build integration

package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/google/uuid"
	"go-microservice-assignment/app/models"
	"testing"
	"time"
)

func TestRedisOptimisticLocking(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "test",
		DB:       0, // use default DB
	})
	db := NewDB(rdb, nil, time.Duration(1)*time.Minute)

	dummyPerson1 := models.Person{
		Id: uuid.New().String(),
		Name: "Test123",
		Address: "Berlin 123",
		DateOfBirth: models.JSONDate(time.Date(1981, time.November, 29, 0, 0, 0, 0, time.UTC)),
	}

	err := db.CreatePerson(ctx, &dummyPerson1)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	dummyPerson1.Name = "Person1"

	// same person but different name to have two separate operation calls
	dummyPerson2 := models.Person{
		Id: dummyPerson1.Id,
		Name: "Person2",
		Address: "Berlin 123",
	}

	updateChan1 := make(chan error)
	updateChan2 := make(chan error)

	trans_failed := false
	for i:=0; i<5; i++ {
		go updatePersonOptimistic1(db, ctx, &dummyPerson1, updateChan1)
		go updatePersonOptimistic2(db, ctx, &dummyPerson2, updateChan2)
		err1, err2 := <-updateChan1, <-updateChan2

		t.Log(err1)
		t.Log(err2)
		if err1 != err2 {
			trans_failed = true
			t.Log("Transaction failed for one of the updates")
			break
		}
	}

	if trans_failed == false {
		t.Fail()
	}
}

func updatePersonOptimistic1(db RedisDB, ctx context.Context, person *models.Person, updateChan1 chan error) {
	_, err := db.UpdatePersonOptimistic(ctx, person)
	updateChan1 <- err
}

func updatePersonOptimistic2(db RedisDB, ctx context.Context, person *models.Person, updateChan2 chan error) {
	_, err := db.UpdatePersonOptimistic(ctx, person)
	updateChan2 <- err
}

func TestRedisPessimisticLocking(t *testing.T) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "test",
		DB:       0, // use default DB
	})

	pool := goredis.NewPool(rdb)
	rs := redsync.New(pool)
	mutex := rs.NewMutex("update-person-lock")

	db := NewDB(rdb, mutex, time.Duration(1)*time.Minute)

	dummyPerson1 := models.Person{
		Id: uuid.New().String(),
		Name: "Test123",
		Address: "Berlin 123",
		DateOfBirth: models.JSONDate(time.Date(1981, time.November, 29, 0, 0, 0, 0, time.UTC)),
	}

	err := db.CreatePerson(ctx, &dummyPerson1)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	dummyPerson1.Name = "Person1"

	// same person but different name to have two separate operation calls
	dummyPerson2 := models.Person{
		Id: dummyPerson1.Id,
		Name: "Person2",
		Address: "Berlin 123",
	}

	updateChanP1 := make(chan error)
	updateChanP2 := make(chan error)

	all_trans_ok := true
	for i:=0; i<5; i++ {
		go updatePersonPessimistic1(db, ctx, &dummyPerson1, updateChanP1)
		go updatePersonPessimistic2(db, ctx, &dummyPerson2, updateChanP2)
		err1, err2 := <-updateChanP1, <-updateChanP2

		t.Log(err1)
		t.Log(err2)
		if err1 != err2 {
			all_trans_ok = false
			t.Log("Transaction failed for one of the updates")
			break
		}
	}

	if all_trans_ok == false {
		t.Fail()
	}
}

func updatePersonPessimistic1(db RedisDB, ctx context.Context, person *models.Person, updateChanP1 chan error) {
	_, err := db.UpdatePersonPessimistic(ctx, person)
	updateChanP1 <- err
}

func updatePersonPessimistic2(db RedisDB, ctx context.Context, person *models.Person, updateChanP2 chan error) {
	_, err := db.UpdatePersonPessimistic(ctx, person)
	updateChanP2 <- err
}

