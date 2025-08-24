package models

import "time"

type Status string

const (
	StatusBorrowed Status = "borrowed"
	StatusReturned Status = "returned"
)

// this borrow record ties a specific user to a specific book at a given time
type Borrow struct {
	BaseModel
	UserID     uint       `gorm:"not null; index" json:"user_id"` // FK to users, enforces WHO borrowed a book.
	BookID     uint       `gorm:"not null; index" json:"book_id"` // FK to books, enforces WHICH book was borrowed.
	BorrowDate time.Time  `gorm:"not null" json:"borrow_date"`
	DueDate    time.Time  `gorm:"not null" json:"due_date"`
	ReturnDate *time.Time `json:"return_date"` // nullable, since the book may not be returned
	Status     Status     `gorm:"type:enum('borrowed', 'returned'); default:'borrowed'" json:"status"`

	// RELATIONSHIPS
	// the first User is the field name, allows me to access sth like borrow.User
	// the second user refers to the User struct
	User User `gorm:"foreignkey:UserID" json:"user"`
	Book Book `gorm:"foreignkey:BookID" json:"book"`
}

func (Borrow) TableName() string {
	return "borrows"
}
