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

type BorrowService struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewBorrowService(db *gorm.DB, redisClient *redis.Client) BorrowService {

	return BorrowService{DB: db, Redis: redisClient}
}

// service to get all books for member
func (b *BorrowService) GetPaginatedMemberBooksService(r *http.Request) (*common.Pagination, error) {

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
	result := b.DB.Scopes(pagination.Paginate()).Order("created_at desc").Find(&books)

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

// service to handle borrow logic

func (b *BorrowService) BorrowBookService(userId uint, payload *requests.BorrowRequest) (*models.Borrow, error) {

	// validate the user id if present
	var user models.User
	result := b.DB.First(&user, userId)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Default().Println("error getting user", result.Error)
		return nil, errors.New("error getting user")
	}

	// validate the book id if present
	var book models.Book
	result = b.DB.First(&book, payload.BookID)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Default().Println("error getting book", result.Error)
		return nil, errors.New("error getting book")
	}

	// check if book stock is available
	if book.Stock < 1 {
		return nil, errors.New("book is not available")
	}

	var existingBorrow models.Borrow

	if err := b.DB.Where("user_id = ? AND book_id = ? AND status = ?", userId, payload.BookID, models.StatusBorrowed).First(&existingBorrow).Error; err == nil {

		return nil, errors.New("book is already borrowed")
	}

	const maxBooksAllowed = 5
	var activeBorrows int64
	b.DB.Model(&models.Borrow{}).Where("user_id = ? AND status = ?", userId, models.StatusBorrowed).Count(&activeBorrows)

	if activeBorrows >= maxBooksAllowed {
		return nil, errors.New("user has already borrowed 5 books")
	}

	// if all the conditions are passed then create a new borrow record

	borrow := &models.Borrow{
		UserID:     userId,
		BookID:     payload.BookID,
		BorrowDate: time.Now(),
		DueDate:    payload.DueDate,
		ReturnDate: nil,
		Status:     models.StatusBorrowed,
	}

	result = b.DB.Create(borrow)

	if result.Error != nil {
		log.Default().Println("error creating borrow", result.Error)
		return nil, errors.New("error creating borrow record")
	}

	if err := b.DB.Preload("User").Preload("Book").First(&borrow, borrow.ID).Error; err != nil {
		log.Default().Println("error getting borrow record", err)
		// return nil, errors.New("error getting borrow record")
	}
	// decrement book stock
	book.Stock -= 1

	if err := b.DB.Save(&book).Error; err != nil {
		log.Default().Println("error decrementing book stock", err)
		return nil, errors.New("error decrementing book stock")
	}

	return borrow, nil

}
