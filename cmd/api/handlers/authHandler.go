package handlers

import (
	// "errors"
	"fmt"
	"log"

	"github.com/Wanjie-Ryan/LMS/cmd/api/requests"
	"github.com/Wanjie-Ryan/LMS/cmd/api/services"
	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/Wanjie-Ryan/LMS/internal/models"
	"github.com/labstack/echo/v4"
)

// Handler to register a new user
// RegisterUserHandler godoc
// @Summary      Register user
// @Description  Creates a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      requests.RegisterRequest  true  "Register payload"
// @Success      201  {object}  common.JsonSuccessResponse
// @Failure      400  {object}  common.JsonErrorResponse  "Invalid payload / user exists"
// @Failure 422 {object} common.JsonFailedValidationResponse
// @Failure      500  {object}  common.JsonErrorResponse  "Server error"
// @Router       /auth/register [post]
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

	return common.SendCreatedResponse(c, "User registered successfully", savedUser)
}

// Handler to Login a user
// LoginHandler godoc
// @Summary      Login
// @Description  Logs in a user and returns access & refresh tokens
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        payload  body      requests.LoginRequest  true  "Login payload"
// @Success      200  {object}  common.JsonSuccessResponse  "accessToken, refreshToken, user"
// @Failure      400  {object} common.JsonErrorResponse  "Invalid credentials"
// @Failure      404  {object}  common.JsonErrorResponse  "User not found"
// @Failure 422 {object} common.JsonFailedValidationResponse
// @Failure      500  {object}  common.JsonErrorResponse  "Server error"
// @Router       /auth/login [post]
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
// ProfileHandler godoc
// @Summary      Get profile
// @Description  Returns the current user's profile
// @Tags         Profile
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  common.JsonSuccessResponse  "user"
// @Failure      401  {object}  common.JsonErrorResponse  "Not authorized"
// @Failure      404  {object}  common.JsonErrorResponse  "User not found"
// @Failure      500  {object} common.JsonErrorResponse  "Server error"
// @Router       /profile/lookup [get]
func (h *Handler) ProfileHandler(c echo.Context) error {

	// get the context of the current logged in user
	user, ok := c.Get("user").(*models.User)
	if !ok {
		return common.SendUnauthorizedResponse(c, "Not authorized")
	}

	authService := services.NewAuthService(h.DB, h.Redis)

	// getting the user by the id
	user, err := authService.ProfileLookupService(user.ID)
	if err != nil {
		log.Default().Println("error getting user", err.Error())
		if err.Error() == "error getting user" {
			return common.SendInternalServerError(c, "error getting user")
		}
	}

	if user == nil {
		return common.SendNotFoundResponse(c, "User does not exist")
	}

	return common.SendSuccessResponse(c, "User profile", user)
}
