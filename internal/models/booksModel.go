package models

type Book struct {
	BaseModel
	Title       string  `gorm:"type:varchar(200); not null; uniqueIndex" json:"title"`
	Author      string  `gorm:"type:varchar(200); not null" json:"author"`
	Description *string `gorm:"type:text" json:"description"`
	Stock       uint    `gorm:"type:int unsigned; not null" json:"stock"`
}

func (Book) TableName() string {
	return "books"
}
