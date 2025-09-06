package models

type Book struct {
	BaseModel
	Title       string  `gorm:"type:varchar(200); not null; uniqueIndex" json:"title"`
	Author      string  `gorm:"type:varchar(200); not null" json:"author"`
	Description *string `gorm:"type:text" json:"description"`
	Stock       uint    `gorm:"type:int unsigned; not null" json:"stock"`
	UserID      uint    `gorm:"not null; index" json:"user_id"`

	User User `gorm:"foreignkey:UserID" json:"user"`
}

func (Book) TableName() string {
	return "books"
}
