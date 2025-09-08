package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

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

func (b *BooksService) CreateBooksService(payload *requests.BookRequest, userId uint) (*models.Book, error) {

	savedBook := &models.Book{
		Title:       payload.Title,
		Author:      payload.Author,
		Description: payload.Description,
		Stock:       payload.Stock,
		UserID:      userId,
	}

	// saving the data to DB
	result := b.DB.Create(&savedBook)
	if result.Error != nil {
		log.Default().Println("error when creating book", result.Error)
		return nil, errors.New("error creating book")
	}

	// saving the data to Redis, first refetch the data with the preloaded user
	var preloadedBook models.Book
	b.DB.Preload("User").First(&preloadedBook, savedBook.ID)
	booksJson, err := json.Marshal(preloadedBook)
	if err != nil {
		fmt.Println("error marshalling book struct to json", err)
	} else {
		err = b.Redis.Set(common.Ctx, fmt.Sprintf("book:%d", savedBook.ID), booksJson, time.Minute*5).Err()

		if err != nil {
			log.Default().Println("error saving book to redis", err)
		} else {
			log.Default().Println("book saved to redis successfully")
		}
	}

	return &preloadedBook, nil
}

// get all paginated books

func (b *BooksService) GetPaginatedBooksService(r *http.Request) (*common.Pagination, error) {

	q := r.URL.Query()
	page := q.Get("page")
	limit := q.Get("limit")
	cacheKey := fmt.Sprintf("books:page:%s:limit:%s", page, limit)

	val, err := b.Redis.Get(common.Ctx, cacheKey).Result()
	if err == nil && val != "" {
		var paginated common.Pagination
		if jsonErr := json.Unmarshal([]byte(val), &paginated); jsonErr == nil {
			log.Default().Println("paginated books fetched from redis successfully")
			return &paginated, nil
		}
	}

	pagination := common.NewPagination(&models.Book{}, r, b.DB)

	var books []models.Book
	result := b.DB.Preload("User").Scopes(pagination.Paginate()).Order("created_at desc").Find(&books)
	if result.Error != nil {
		log.Default().Println("error getting books", result.Error)
		return nil, errors.New("error getting books")
	}

	pagination.Data = books

	booksJson, err := json.Marshal(pagination)
	if err != nil {
		fmt.Println("error marshalling book struct to json in redis", err)
	} else {
		// err = b.Redis.Set(common.Ctx, cacheKey, booksJson, 0).Err()
		err = b.Redis.Set(common.Ctx, cacheKey, booksJson, time.Minute*2).Err()
	}

	return pagination, nil

}

// update books

func (b *BooksService) UpdateBooksService(payload *requests.UpdateBookRequest, userId uint) (*models.Book, error) {

	cacheKey := fmt.Sprintf("book:%d", payload.ID)

	var books models.Book
	result := b.DB.First(&books, payload.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Default().Println("error getting book", result.Error)
		return nil, errors.New("error getting book")
	}

	if payload.Title != nil {
		books.Title = *payload.Title
	}

	if payload.Author != nil {
		books.Author = *payload.Author
	}

	if payload.Description != nil {
		books.Description = payload.Description
	}

	if payload.Stock != nil {
		books.Stock = *payload.Stock
	}

	books.UserID = userId
	books.UpdatedAt = time.Now()

	// updatedResult := b.DB.Preload("User").Save(&books)
	updatedResult := b.DB.Save(&books)
	if updatedResult.Error != nil {
		log.Default().Println("error updating book in db", updatedResult.Error)
		return nil, errors.New("error updating book in db")
	}

	// saving to redis
	var preloadedBook models.Book
	b.DB.Preload("User").First(&preloadedBook, payload.ID)
	bookJson, err := json.Marshal(preloadedBook)
	if err != nil {
		log.Default().Println("error marshalling book struct to json", err)
	} else {
		err = b.Redis.Set(common.Ctx, cacheKey, bookJson, time.Minute*5).Err()
		if err != nil {
			log.Default().Println("error updating book to redis", err)
		} else {
			log.Default().Println("book updated to redis successfully")
		}

	}

	return &books, nil

}

// get single book
func (b *BooksService) GetSingleBookService(id uint) (*models.Book, error) {

	var book models.Book
	cachekey := fmt.Sprintf("book:%d", id)

	// first check redis

	val, err := b.Redis.Get(common.Ctx, cachekey).Result()
	if err == nil && val != "" {
		var cachedBook models.Book
		if jsonErr := json.Unmarshal([]byte(val), &cachedBook); jsonErr == nil {
			log.Default().Println("book fetched from redis successfully")
			return &cachedBook, nil
		} else {
			log.Default().Println("error unmarshalling book from redis", jsonErr)
		}

	}

	result := b.DB.Preload("User").First(&book, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Default().Println("error getting book", result.Error)
		return nil, errors.New("error getting book")
	}

	bookJson, err := json.Marshal(book)
	if err != nil {
		log.Default().Println("error marshalling book struct to json", err)
	} else {
		err = b.Redis.Set(common.Ctx, cachekey, bookJson, time.Minute*5).Err()
		if err != nil {
			log.Default().Println("error saving book to redis", err)
		} else {
			log.Default().Println("book saved to redis successfully")
		}
	}

	return &book, nil

}

// delete a book

func (b *BooksService) DeleteBooksService(id uint) error {
	cachekey := fmt.Sprintf("book:%d", id)

	result := b.DB.Delete(&models.Book{}, id)
	if result.Error != nil {

		log.Default().Println("error deleting book", result.Error)
		return errors.New("error deleting book")
	}
	if result.RowsAffected == 0 {
		return errors.New("book not found")
	}

	err := b.Redis.Del(common.Ctx, cachekey).Err()

	if err != nil {
		log.Default().Println("error deleting book from redis", err)
	} else {
		log.Default().Println("book deleted from redis successfully")
	}
	return nil

}
