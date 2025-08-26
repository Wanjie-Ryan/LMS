package services

import (
	"errors"
	"fmt"

	"github.com/Wanjie-Ryan/LMS/cmd/api/requests"
	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/Wanjie-Ryan/LMS/internal/models"
	"gorm.io/gorm"
)

type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {

	return AuthService{DB: db}
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
