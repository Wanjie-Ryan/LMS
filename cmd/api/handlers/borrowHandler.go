package handlers

import (
	// "log"
	"log"

	"github.com/labstack/echo/v4"

	"github.com/Wanjie-Ryan/LMS/cmd/api/requests"
	"github.com/Wanjie-Ryan/LMS/cmd/api/services"
	"github.com/Wanjie-Ryan/LMS/common"
	"github.com/Wanjie-Ryan/LMS/internal/models"
)

// handler for member to get paginated books
// GetPaginatedBooksHandler godoc
// @Summary      Get paginated books
// @Description  Returns a paginated list of books for members
// @Tags         Borrow Books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page  query     string  false  "Page number"  default(1)
// @Param        limit  query     string  false  "Limit"  default(10)
// @Success      200  {object}  common.JsonSuccessResponse
// @Failure      401  {object}  common.JsonErrorResponse  "Not authorized"
// @NotFound     404  {object}  common.JsonErrorResponse  "Not found"
// @Failure      500  {object}  common.JsonErrorResponse  "Server error"
// @Router       /books/getAll [get]
func (h *Handler) GetMemberPaginatedBooksHandler(c echo.Context) error {
	user, ok := c.Get("user").(*models.User)
	// log.Default().Println("admin email is checking", user.Email)
	if !ok {
		return common.SendUnauthorizedResponse(c, "Not authorized")
	}

	if user.Role != "member" {
		return common.SendForbiddenResponse(c, "Not allowed to perform this action")
	}

	borrowService := services.NewBorrowService(h.DB, h.Redis)

	books, err := borrowService.GetPaginatedMemberBooksService(c.Request())
	if err != nil {
		if err.Error() == "error getting books" {
			return common.SendInternalServerError(c, "Error getting books")
		}
	}

	return common.SendSuccessResponse(c, "Books retrieved successfully", books)
}

// handler for user to borrow book
// BorrowHandler godoc
// @Summary      Member can borrow books
// @Description  Member is able to borrow books
// @Tags         Borrow Books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      requests.BorrowRequest  true  "Borrow Book payload"
// @Success      201  {object}  common.JsonSuccessResponse
// @Failure      401  {object}  common.JsonErrorResponse  "Not authorized"
// @Failure      403  {object}  common.JsonErrorResponse  "Forbidden"
// @NotFound     404  {object}  common.JsonErrorResponse  "Not found"
// @Failure      500  {object}  common.JsonErrorResponse  "Server error"
// @Router       /borrow [post]
func (h *Handler) BorrowHandler(c echo.Context) error {

	user, ok := c.Get("user").(*models.User)
	if !ok {
		return common.SendUnauthorizedResponse(c, "Not authorized")
	}

	if user.Role != "member" {

		return common.SendForbiddenResponse(c, "Not allowed to perform this action")
	}

	payload := new(requests.BorrowRequest)

	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, "Invalid request payload")
	}

	validationErr := h.ValidateBodyRequest(c, payload)
	if validationErr != nil {
		return common.SendFailedValidationResponse(c, validationErr)
	}

	borrowService := services.NewBorrowService(h.DB, h.Redis)

	borrow, err := borrowService.BorrowBookService(user.ID, payload)

	if err != nil {
		if err.Error() == "error borrowing book" {
			return common.SendInternalServerError(c, "Error borrowing book")
		} else if err.Error() == "book is already borrowed" {
			return common.SendForbiddenResponse(c, "Book already borrowed")
		} else if err.Error() == "book is not available" {
			return common.SendForbiddenResponse(c, "Book is not available")
		} else if err.Error() == "error getting book" {
			return common.SendInternalServerError(c, "Error getting book")
		} else if err.Error() == "error getting user" {
			return common.SendInternalServerError(c, "Error getting member")
		} else if err.Error() == "user has already borrowed 5 books" {
			return common.SendForbiddenResponse(c, "User has already borrowed 5 books")
		} else if err.Error() == "error creating borrow record" {
			return common.SendInternalServerError(c, "Error creating borrow record")
		} else {
			log.Default().Println("error borrowing book", err)
			return common.SendInternalServerError(c, "Error borrowing book")
		}
	}
	if borrow == nil {
		return common.SendNotFoundResponse(c, "Book not found")
	}

	return common.SendCreatedResponse(c, "Book borrowed successfully", borrow)
}

// handler for user to return book
// ReturnBookHandler godoc
// @Summary      Member can Return books
// @Description  Member is able to Return books
// @Tags         Return Books
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        payload  body      requests.ReturnRequest  true  "Return Book payload"
// @Success      200  {object}  common.JsonSuccessResponse
// @Failure      401  {object}  common.JsonErrorResponse  "Not authorized"
// @Failure      403  {object}  common.JsonErrorResponse  "Forbidden"
// @NotFound     404  {object}  common.JsonErrorResponse  "Not found"
// @Failure      500  {object}  common.JsonErrorResponse  "Server error"
// @Router       /return [post]
func (h *Handler) ReturnBookHandler(c echo.Context) error {

	user, ok := c.Get("user").(*models.User)

	if !ok {
		return common.SendUnauthorizedResponse(c, "Not authorized")
	}

	if user.Role != "member" {
		return common.SendForbiddenResponse(c, "Not allowed to perform this action")
	}

	payload := new(requests.ReturnRequest)

	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		return common.SendBadRequestResponse(c, "Invalid request payload")
	}

	validationErr := h.ValidateBodyRequest(c, payload)

	if validationErr != nil {
		return common.SendFailedValidationResponse(c, validationErr)
	}

	borrowService := services.NewBorrowService(h.DB, h.Redis)

	borrow, err := borrowService.ReturnBookService(user.ID, payload)

	if err != nil {
		if err.Error() == "error returning book" {
			return common.SendInternalServerError(c, "Error returning book")
		} else if err.Error() == "error getting book" {
			return common.SendInternalServerError(c, "Error getting book")
		} else if err.Error() == "error getting user" {
			return common.SendInternalServerError(c, "Error getting member")
		} else if err.Error() == "error getting borrow record" {
			return common.SendInternalServerError(c, "Error getting borrow record")
		} else {
			log.Default().Println("error returning book", err)
			return common.SendInternalServerError(c, "Error returning book")
		}
	}

	if borrow == nil {
		return common.SendNotFoundResponse(c, "Book not found")
	}

	return common.SendSuccessResponse(c, "Book returned successfully", borrow)
}
