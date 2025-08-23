package common

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// making the connection to my DB, this will return a gorm db connection or an error

func ConnectionDB() (*gorm.DB, error) {

	// loading the .env files
	err := godotenv.Load()
	if err != nil {
		panic("Error loading the .env file")
		// the return will be unreachable because of the panic
		// return nil, err
	}

	// getting the database connections from the env files
	host := os.Getenv("DB_HOST")
	db_name := os.Getenv("DB_NAME")
	db_username := os.Getenv("DB_USERNAME")
	db_password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", db_username, db_password, host, db_name)
	fmt.Println(dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		fmt.Println("Failed to connect to database", err)
		// panic("Failed to connect to database", err)
		return nil, err
	}

	log.Default().Println("Connected to the database")

	return db, nil

}
