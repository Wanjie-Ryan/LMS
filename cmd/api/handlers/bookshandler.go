package handlers

import (
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/Wanjie-Ryan/LMS/cmd/api/requests"
	"github.com/Wanjie-Ryan/LMS/cmd/api/services"
	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/Wanjie-Ryan/LMS/internal/models"
)

// Handler to create a book
// CreateBookHandler godoc
// @Summary      Create book
// @Description  Allows admin users to create a new book entry.
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        payload  body      requests.BookRequest  true  "Book payload"
// @Success      201  {object}  common.JsonSuccessResponse
// @Failure      400  {object}  common.JsonErrorResponse  "Invalid payload"
// @Failure      401  {object}  common.JsonErrorResponse  "Not authorized"
// @Failure      403  {object}  common.JsonErrorResponse  "Forbidden"
// @Failure 422 {object} common.JsonFailedValidationResponse
// @Failure      500  {object}  common.JsonErrorResponse  "Server error"
// @Router       /books/create [post]
func (h *Handler) CreateBookHandler(c echo.Context) error {

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return common.SendUnauthorizedResponse(c, "Not authorized")
	}

	if user.Role != "admin" {
		return common.SendForbiddenResponse(c, "Not allowed to perform this action")
	}

	payload := new(requests.BookRequest)

	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		fmt.Println("error binding book payload", err.Error())
		return common.SendBadRequestResponse(c, "Invalid book creation payload")
	}

	validationErr := h.ValidateBodyRequest(c, *payload)

	if validationErr != nil {
		return common.SendFailedValidationResponse(c, validationErr)
	}

	booksService := services.NewBookService(h.DB, h.Redis)

	book, err := booksService.CreateBooksService(payload, user.ID)
	if err != nil {
		if err.Error() == "error creating book" {
			return common.SendInternalServerError(c, "Error creating book")
		}
	}

	return common.SendCreatedResponse(c, "Book created successfully", book)
}

// Handler to get paginated books
// GetPaginatedBooksHandler godoc
// @Summary      Get paginated books
// @Description  Returns a paginated list of books
// @Tags         Books
// @Accept       json
// @Produce      json
// @Param        page  query     string  false  "Page number"  default(1)
// @Param        limit  query     string  false  "Limit"  default(10)
// @Success      200  {object}  common.JsonSuccessResponse
// @Failure      401  {object}  common.JsonErrorResponse  "Not authorized"
// @NotFound     404  {object}  common.JsonErrorResponse  "Not found"
// @Failure      500  {object}  common.JsonErrorResponse  "Server error"
// @Router       /books [get]
