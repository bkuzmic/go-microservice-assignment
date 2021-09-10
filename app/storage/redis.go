package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-microservice-assignment/app/models"
	"time"
)

type RedisDB interface {
	CreatePerson(ctx *context.Context, p *models.Person) error
	GetPerson(ctx *context.Context, id string) (*models.Person, error)
	UpdatePersonOptimistic(ctx *context.Context, p *models.Person) error
	UpdatePersonPessimistic(ctx *context.Context, p *models.Person) error
}

type DB struct {
	Client *redis.Client
	ExpireTimeInMinutes time.Duration
}