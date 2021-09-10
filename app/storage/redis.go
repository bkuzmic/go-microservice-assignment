package storage

import (
	"github.com/go-redis/redis/v8"
	"go-microservice-assignment/app/models"
)

type RedisDB interface {
	CreatePerson(p *models.Person) error
	GetPerson(id string) (*models.Person, error)
	UpdatePersonOptimistic(p *models.Person) error
	UpdatePersonPessimistic(p *models.Person) error
}

type DB struct {
	Client *redis.Client
}