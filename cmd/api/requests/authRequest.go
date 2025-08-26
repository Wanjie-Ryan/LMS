package requests

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

// REGISTRATION DTO
type RegisterRequest struct {
	Firstname string `json:"firstname" validate:"required"`
	Lastname  string `json:"lastname" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=4"`
	Role      Role   `json:"role" validate:"omitempty,oneof=admin member"`
	//oneof means that the field must be one of the values in the list
	// in my case, if a role is provided, it must be either admin or member, if nothing is provided, it defaults to member from the model
}
