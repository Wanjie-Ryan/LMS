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
	Title       *string `json:"title,omitempty"`
	Author      *string `json:"author,omitempty"`
	Description *string `json:"description,omitempty"`
	Stock       *uint   `json:"stock,omitempty"`
}

// struct to delete and get by id
type GetDelBookRequest struct {
	ID uint `json:"id" validate:"required"`
}
