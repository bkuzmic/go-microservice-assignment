package storage

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go-microservice-assignment/app/models"
	"time"
)

type db struct {
	client *redis.Client
	expireTimeInMinutes time.Duration
}

type RedisDB interface {
	CreatePerson(ctx context.Context, p *models.Person) error
	GetPerson(ctx context.Context, id string) (*models.Person, error)
	UpdatePersonOptimistic(ctx context.Context, p *models.Person) (*models.Person, error)
	//UpdatePersonPessimistic(ctx *context.Context, p *models.Person) (*models.Person, error)
}

func NewDB(client *redis.Client, expireTimeInMinutes time.Duration) RedisDB {
	return &db{client, expireTimeInMinutes}
}

func (d *db) CreatePerson(ctx context.Context, p *models.Person) error {
	expireKey := getExpireKey(p.Id)
	created := time.Now()

	trans := d.client.TxPipeline()
	// insert person with person.Id as key
	trans.Set(ctx, p.Id, p, 0)
	// also insert key with updated date and expiration
	trans.Set(ctx, expireKey, created, d.expireTimeInMinutes)
	_, err := trans.Exec(ctx)

	return err
}

func (d *db) GetPerson(ctx context.Context, id string) (*models.Person, error) {
	var person models.Person
	res, err := d.client.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(res), &person)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (d *db) UpdatePersonOptimistic(ctx context.Context, p *models.Person) (*models.Person, error) {
	var modifiedPerson *models.Person

	err := d.client.Watch(ctx, func(tx *redis.Tx) error {
		personString, err := tx.Get(ctx, p.Id).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		err = json.Unmarshal([]byte(personString), &modifiedPerson)
		if err != nil {
			return err
		}

		// update person's data
		if p.Name != "" {
			modifiedPerson.Name = p.Name
		}
		if p.Address != "" {
			modifiedPerson.Address = p.Address
		}
		dateOfBirth, err := p.DateOfBirth.MarshalText()
		if err == nil && dateOfBirth != "01/01/0001" {
			modifiedPerson.DateOfBirth = p.DateOfBirth
		}

		expireKey := getExpireKey(modifiedPerson.Id)
		updated := time.Now()

		trans := tx.TxPipeline()
		// insert person with person.Id as key
		trans.Set(ctx, modifiedPerson.Id, modifiedPerson, 0)
		// also insert key with updated date and expiration
		trans.Set(ctx, expireKey, updated, d.expireTimeInMinutes)
		_, err = trans.Exec(ctx)

		return err
	}, p.Id)

	return modifiedPerson, err
}

func getExpireKey(id string) string {
	return id + "_expire"
}