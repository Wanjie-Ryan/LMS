package handlers

import (
	"fmt"

	"github.com/Wanjie-Ryan/LMS/cmd/api/requests"
	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/labstack/echo/v4"
)

// Handler to register a new user
func (h *Handler) RegisterUserHandler(c echo.Context) error{

	payload := new(requests.RegisterRequest)

	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err !=nil{
		fmt.Println("error binding body", err.Error())
		return common.SendBadRequestResponse(c, "Invalid Payload")
	
	}

	validationErr := h.ValidateBodyRequest(c, payload)

	return nil
}