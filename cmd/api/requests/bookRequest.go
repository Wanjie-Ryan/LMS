package requests

type BookRequest struct {
	Title       string  `json:"title" validate:"required"`
	Author      string  `json:"author" validate:"required"`
	Description *string `json:"description"`
	Stock       uint    `json:"stock" validate:"required"`
}
