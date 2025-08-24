package handlers

// the handler will be instantiated in main.go with the actual db instance, and the handler will be reused globally.
import "gorm.io/gorm"

type Handler struct {
	DB *gorm.DB
}
