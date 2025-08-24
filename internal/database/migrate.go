package main

import (
	"log"

	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/Wanjie-Ryan/LMS/internal/models"
)

func main() {

	db, err := common.ConnectionDB()
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Book{})

	if err != nil {
		panic(err)
	}
	log.Default().Println("Database Migrated successfully")

}
