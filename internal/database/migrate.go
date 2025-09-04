package database

import (
	"log"

	"gorm.io/gorm"

	// "github.com/Wanjie-Ryan/LMS/common"
	"github.com/Wanjie-Ryan/LMS/internal/models"
)

func Migrate(db *gorm.DB) error {

	// db, err := common.ConnectionDB()
	// if err != nil {
	// 	panic(err)
	// }

	err := db.AutoMigrate(&models.User{}, &models.Book{}, &models.Borrow{})

	if err != nil {
		panic(err)
	}
	log.Default().Println("Database Migrated successfully")
	return nil

}
