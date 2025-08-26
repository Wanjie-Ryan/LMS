package common

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type JsonSuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type JsonFailedValidationResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Errors  []*ValidationError `json:"errors"`
}

type JsonErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// sending a response if the request is successful
func SendSuccessResponse(c echo.Context, message string, data interface{}) error {

	return c.JSON(http.StatusOK, JsonSuccessResponse{Success: true, Message: message, Data: data})
}

// sending a response if validation failed
func SendFailedValidationResponse(c echo.Context, errors []*ValidationError) error {

	return c.JSON(http.StatusUnprocessableEntity, JsonFailedValidationResponse{Success: false, Message: "Validation Failed", Errors: errors})
}

func SendErrorResponse(c echo.Context, statusCode int, message string) error {
	return c.JSON(statusCode, JsonErrorResponse{Success: false, Message: message})
}

func SendUnauthorizedResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, http.StatusUnauthorized, message)
}

func SendBadRequestResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, http.StatusBadRequest, message)
}

func SendNotFoundResponse(c echo.Context, message string) error {
	return SendErrorResponse(c, http.StatusNotFound, message)
}
func SendInternalServerError(c echo.Context, message string) error {
	return SendErrorResponse(c, http.StatusInternalServerError, message)
}
