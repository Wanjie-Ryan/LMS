package models

// creates a new type called Role whose underlying type is string.
// Role is not just a plain string anymore - it's a custom type.
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

// RoleAdmin is a constant of type Role with the value "admin".
// RoleMember is a constant of type Role with the value "member".
// this gives a stronger type safety. Instead of letting anyone assign just any string to Role, we're forcing it to be either "admin" or "member".

type User struct {
	BaseModel
	Firstname string `gorm:"type:varchar(200); not null" json:"firstname"`
	Lastname  string `gorm:"type:varchar(200); not null" json:"lastname"`
	Email     string `gorm:"type:varchar(100); not null; uniqueIndex" json:"email"`
	Password  string `gorm:"type:varchar(200); not null" json:"-"`
	Role      Role   `gorm:"type:enum('admin', 'member'); default:'member'" json:"role"`
}

func (User) TableName() string {
	return "users"
}
