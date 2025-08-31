package services

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type BookService struct {
	DB    *gorm.DB
	Redis *redis.Client
}
