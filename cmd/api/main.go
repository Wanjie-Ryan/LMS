package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal("Error while loading the .env file", err)
	}

	// getting the db connection from the common package
	db, err := common.ConnectionDB()
	if err != nil {
		e.Logger.Fatal("Error while connecting to the database", err)
	}

	fmt.Println("Connected to the database", db)

	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is up and running")

	})

	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf(":%s", port)
	fmt.Println("Server is running on port:", port)
	e.Logger.Fatal(e.Start(appAddress))

}
