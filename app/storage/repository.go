package storage

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go-microservice-assignment/app/models"
	"time"
)

func (d *DB) CreatePerson(ctx context.Context, p *models.Person) error {
	expireKey := GetExpireKey(p.Id)
	created := time.Now()

	trans := d.Client.TxPipeline()
	// insert person with person.Id as key
	trans.Set(ctx, p.Id, p, 0)
	// also insert key with updated date and expiration
	trans.Set(ctx, expireKey, created, d.ExpireTimeInMinutes)
	_, err := trans.Exec(ctx)

	return err
}

func (d *DB) GetPerson(ctx context.Context, id string) (*models.Person, error) {
	var person models.Person
	res, err := d.Client.Get(ctx, id).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(res), &person)
	if err != nil {
		return nil, err
	}
	return &person, nil
}

func (d *DB) UpdatePersonOptimistic(ctx context.Context, p *models.Person) (models.Person, error) {
	var modifiedPerson *models.Person

	err := d.Client.Watch(ctx, func(tx *redis.Tx) error {
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

		expireKey := GetExpireKey(modifiedPerson.Id)
		updated := time.Now()

		trans := tx.TxPipeline()
		// insert person with person.Id as key
		trans.Set(ctx, modifiedPerson.Id, modifiedPerson, 0)
		// also insert key with updated date and expiration
		trans.Set(ctx, expireKey, updated, d.ExpireTimeInMinutes)
		_, err = trans.Exec(ctx)

		return err
	}, p.Id)

	return *modifiedPerson, err
}

func GetExpireKey(id string) string {
	return id + "_expire"
}