package requests

type BookRequest struct {
	Title       string  `json:"title" validate:"required"`
	Author      string  `json:"author" validate:"required"`
	Description *string `json:"description"`
	Stock       uint    `json:"stock" validate:"required"`
}

// struct to update a book
type UpdateBookRequest struct {
	ID          uint    `json:"id" validate:"required"`
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description *string `json:"description"`
	Stock       uint    `json:"stock"`
}
