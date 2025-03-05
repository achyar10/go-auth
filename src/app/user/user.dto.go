package user

type CreateUserDTO struct {
	Username string  `json:"username" validate:"required,min=3,max=100"`
	Password string  `json:"password" validate:"required,min=8"`
	Fullname *string `json:"fullname"`
	Role     Role    `json:"role" validate:"oneof=admin user"`
	IsActive *bool   `json:"is_active"`
}
