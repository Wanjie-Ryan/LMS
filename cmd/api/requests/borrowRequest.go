package requests

import "time"

type BorrowRequest struct {
	BookID  uint      `json:"book_id" validate:"required"`
	DueDate time.Time `json:"due_date" validate:"required"`
}
