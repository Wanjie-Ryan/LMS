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
		fmt.Println("error binding body", err.Error())
		return common.SendBadRequestResponse(c, "Invalid Payload")

	}

	validationErr := h.ValidateBodyRequest(c, *payload)
	if validationErr != nil {
		return common.SendFailedValidationResponse(c, validationErr)
	}

	authService := services.NewAuthService(h.DB)
	fmt.Println("authService", authService)

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
