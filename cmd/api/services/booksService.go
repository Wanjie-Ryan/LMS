package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/Wanjie-Ryan/LMS/cmd/api/requests"
	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/Wanjie-Ryan/LMS/internal/models"
)

type BooksService struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewBookService(db *gorm.DB, redisClient *redis.Client) BooksService {
	return BooksService{DB: db, Redis: redisClient}
}

// function to create books, and persist them in redis for caching

func (b *BooksService) CreateBooksService(payload *requests.BookRequest) (*models.Book, error) {

	savedBook := &models.Book{
		Title:       payload.Title,
		Author:      payload.Author,
		Description: payload.Description,
		Stock:       payload.Stock,
	}

	// saving the data to DB
	result := b.DB.Create(&savedBook)
	if result.Error != nil {
		log.Default().Println("error when creating book", result.Error)
		return nil, errors.New("error creating book")
	}

	// saving the data to Redis
	booksJson, err := json.Marshal(savedBook)
	if err != nil {
		fmt.Println("error marshalling book struct to json", err)
	} else {
		err = b.Redis.Set(common.Ctx, fmt.Sprintf("book:%d", savedBook.ID), booksJson, 0).Err()

		if err != nil {
			log.Default().Println("error saving book to redis", err)
		} else {
			log.Default().Println("book saved to redis successfully")
		}
	}

	return savedBook, nil
}
