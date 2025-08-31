package main

import (
	"github.com/Wanjie-Ryan/LMS/cmd/api/handlers"
)

func (app *Application) Routes(handler handlers.Handler) {

	apiGroup := app.server.Group("/api/v1")

	// authroutes
	authRoutes := apiGroup.Group("/auth")
	// fmt.Println("authRoutes", authRoutes)
	authRoutes.POST("/register", handler.RegisterUserHandler)
	authRoutes.POST("/login", handler.LoginHandler)

	// profile Look up route, will be in the same handler as register and login though.
	profileRoute := apiGroup.Group("/profile", app.AuthMiddleware.AuthMiddleware)
	profileRoute.GET("/lookup", handler.ProfileHandler)

}
