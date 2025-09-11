package requests

import "time"

type BorrowRequest struct {
	BookID  uint      `json:"book_id" validate:"required"`
	DueDate time.Time `json:"due_date" validate:"required"`
}

type ReturnRequest struct {
	BookID uint `json:"book_id" validate:"required"`
	// ReturnDate time.Time `json:"return_date" validate:"required"`
}
