package storage

import (
	"context"
	"go-microservice-assignment/app/models"
)

func (d *DB) CreatePerson(p *models.Person) error {
	err := d.Client.Set(context.Background(), p.Id, p, 0).Err()
	return err
}