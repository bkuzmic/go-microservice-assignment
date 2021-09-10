package storage

import (
	"context"
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

func GetExpireKey(id string) string {
	return id + "_expire"
}