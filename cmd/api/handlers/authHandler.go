package handlers

import (
	// "errors"
	"fmt"
	"log"

	"github.com/Wanjie-Ryan/LMS/cmd/api/requests"
	"github.com/Wanjie-Ryan/LMS/cmd/api/services"
	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/labstack/echo/v4"
)

// Handler to register a new user
func (h *Handler) RegisterUserHandler(c echo.Context) error {

	payload := new(requests.RegisterRequest)

	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		fmt.Println("error binding registration body", err.Error())
		return common.SendBadRequestResponse(c, "Invalid Payload")

	}

	validationErr := h.ValidateBodyRequest(c, *payload)
	if validationErr != nil {
		return common.SendFailedValidationResponse(c, validationErr)
	}

	authService := services.NewAuthService(h.DB, h.Redis)

	// check if email exists already
	user, err := authService.GetUserByMail(payload.Email)

	if err != nil {
		log.Default().Println("error getting user", err.Error())
		if err.Error() == "error getting user" {
			return common.SendInternalServerError(c, "error getting user")
		}
	}

	// if the user is not empty
	if user != nil {
		return common.SendBadRequestResponse(c, "User already exists")
	}

	savedUser, err := authService.RegisterService(payload)
	if err != nil {
		log.Default().Println("error registering user", err.Error())
		return common.SendInternalServerError(c, "error registering user")

	}

	return common.SendSuccessResponse(c, "User registered successfully", savedUser)
}

// Handler to Login a user
func (h *Handler) LoginHandler(c echo.Context) error {

	payload := new(requests.LoginRequest)

	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		fmt.Println("error binding login payload", err.Error())
		return common.SendBadRequestResponse(c, "Invalid Payload")
	}

	validationErr := h.ValidateBodyRequest(c, *payload)
	if validationErr != nil {
		return common.SendFailedValidationResponse(c, validationErr)
	}

	authService := services.NewAuthService(h.DB, h.Redis)

	// first check the email
	user, err := authService.GetUserByMail(payload.Email)
	if err != nil {
		log.Default().Println("error getting user during login", err.Error())
		if err.Error() == "error getting user" {
			return common.SendInternalServerError(c, "error getting user")
		}
	}

	if user == nil {
		return common.SendNotFoundResponse(c, "User does not exist, register first")
	}

	// if user exists, check the passwords if they match
	isPasswordMatch := common.ComparePasswords(payload.Password, user.Password)
	if !isPasswordMatch {
		return common.SendBadRequestResponse(c, "Password does not match")
	}

	// generating the jwt
	accessToken, refreshToken, err := common.GenerateJWT(*user)

	if err != nil {
		log.Default().Println("error generating jwt", err.Error())
		return common.SendInternalServerError(c, "error generating jwt")
	}

	return common.SendSuccessResponse(c, "User Logged in sucessfully", map[string]interface{}{"accessToken": accessToken, "refreshToken": refreshToken, "user": user})
}

// handler to Get user profile
func (h *Handler) ProfileHandler(c echo.Context) error {

	return nil
}