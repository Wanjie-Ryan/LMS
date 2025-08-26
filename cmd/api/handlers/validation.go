package handlers

// global package for validation of payloads
import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (h *Handler) ValidateBodyRequest(c echo.Context, payload interface{}) []*common.ValidationError {

	var validate *validator.Validate
	validate = validator.New(validator.WithRequiredStructEnabled())
	var errors []*common.ValidationError

	err := validate.Struct(payload)

	validationErrors, ok := err.(validator.ValidationErrors)

	if ok {

		reflected := reflect.ValueOf(payload)

		for _, validationErr := range validationErrors {

			fmt.Println(reflected.Type().FieldByName(validationErr.StructField()))

			field, _ := reflected.Type().FieldByName(validationErr.StructField())

			key := field.Tag.Get("json")

			if key == "" {
				key = strings.ToLower(validationErr.StructField())
			}

			condition := validationErr.Tag()
			param := validationErr.Param()
			errMessage := key + " field is " + condition

			switch condition {
			case "required":
				errMessage = key + " is required"
			case "email":
				errMessage = key + " must be a valid email address"
			case "min":
				errMessage = key + " must be at least " + param + " characters long"
			case "oneof":
				errMessage = key + " must be one of [" + param + "]"
			}

			currentValidationError := &common.ValidationError{
				Error:     errMessage,
				Key:       key,
				Condition: condition,
			}
			errors = append(errors, currentValidationError)
		}

	}

	return errors

}
