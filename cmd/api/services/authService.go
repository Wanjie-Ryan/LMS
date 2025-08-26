package services

import (
	"github.com/Wanjie-Ryan/LMS/cmd/api/requests"
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

	return nil, nil

}
