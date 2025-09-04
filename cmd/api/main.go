// Package main Library Management System API
//
// @contact.name   Wanjie-Ryan
// @contact.url    https://github.com/Wanjie-Ryan
// @contact.email  ryanwanjie1@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @title           LMS API
// @version         1.0
// @description     Library Management System API with JWT auth, Redis cache, and MySQL.
// @BasePath        /api/v1
// @schemes         http
// @host            localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer {token}"

package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"github.com/Wanjie-Ryan/LMS/cmd/api/handlers"
	"github.com/Wanjie-Ryan/LMS/cmd/api/middleware"
	 "github.com/Wanjie-Ryan/LMS/internal/database"
	"github.com/Wanjie-Ryan/LMS/common"
	_ "github.com/Wanjie-Ryan/LMS/docs"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Application struct {
	logger         echo.Logger
	server         *echo.Echo
	handler        handlers.Handler
	AuthMiddleware middleware.AppMiddleware
}

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
	 if err := database.Migrate(db); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }

	// health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Server is up and running")

	})

	// swagger docs endpoint
	e.GET("/docs/*", echoSwagger.WrapHandler)
	redisClient := common.ConnectRedis()

	handler := handlers.Handler{DB: db, Redis: redisClient}
	authMiddleware := middleware.AppMiddleware{DB: db}

	app := &Application{
		logger:         e.Logger,
		server:         e,
		handler:        handler,
		AuthMiddleware: authMiddleware,
	}

	fmt.Println("app instantiation", app)
	app.Routes(handler)

	// fmt.Println("Connected to the database", db)

	port := os.Getenv("APP_PORT")
	appAddress := fmt.Sprintf(":%s", port)
	fmt.Println("Server is running on port:", port)
	e.Logger.Fatal(e.Start(appAddress))

}
