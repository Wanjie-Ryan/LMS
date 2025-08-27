package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Wanjie-Ryan/LMS/cmd/api/requests"
	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/Wanjie-Ryan/LMS/internal/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewAuthService(db *gorm.DB, redisClient *redis.Client) AuthService {

	return AuthService{DB: db, Redis: redisClient}
}

// function to handle Registration Logic
func (a *AuthService) RegisterService(payload *requests.RegisterRequest) (*models.User, error) {

	// hashing password
	hashedPassword, err := common.HashPassword(payload.Password)

	if err != nil {
		fmt.Println("error hashing password", err)
		return nil, errors.New("error hashing password")
	}

	saveUser := &models.User{
		Firstname: payload.Firstname,
		Lastname:  payload.Lastname,
		Email:     payload.Email,
		Password:  hashedPassword,
		Role:      models.Role(payload.Role),
		// the reason for conversion is that the roles in model and request are from different packages
		// the conversion tells GO, take the Requests.Role string and treat it as a models.Role
	}

	result := a.DB.Create(&saveUser)
	if result.Error != nil {
		fmt.Println("error when saving user to DB during registration", result.Error)
		return nil, errors.New("registration Failed")
	}

	// saving to redis
	//converting the saveUser(struct) to Json string
	// redis cannot store Go structs dorectly, it only stores bytes or strings
	userJson, err := json.Marshal(saveUser)
	if err != nil {
		fmt.Println("error marshalling user for redis", err)
	} else {
		// builds a redis key using fmt.Sprintf, to uniquely identify each user
		// fmt.Sprintf formats strings
		// "user:%d" is a placeholder for the user ID eg. user:1
		// 0 tells redis that the stored data to never expire
		err = a.Redis.Set(common.Ctx, fmt.Sprintf("user:%d", saveUser.ID), userJson, 0).Err()
		if err != nil {
			fmt.Println("error saving user to redis", err)
		} else {
			fmt.Println("user saved to redis successfully")
		}

	}

	return saveUser, nil

}

// function to handle Login Logic
func (a *AuthService) LoginService(email string, password string) (*models.User, error) {
	return nil, nil

}

// function to Get user by email
func (a *AuthService) GetUserByMail(email string) (*models.User, error) {

	var user models.User
	result := a.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
			// by doing nil, nil, the caller distinguishes btn the error from the db, and when user actually doesn't exist
		}

		return nil, errors.New("error getting user")
	}

	return &user, nil
}
