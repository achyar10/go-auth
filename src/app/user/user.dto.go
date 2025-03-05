package user

type CreateUserDTO struct {
	Username string  `json:"username" validate:"required,min=3,max=100"`
	Password string  `json:"password" validate:"required,min=8"`
	Fullname *string `json:"fullname"`
	Role     Role    `json:"role" validate:"oneof=admin user"`
	IsActive *bool   `json:"is_active"`
}

type ListUserQueryDTO struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	SortBy string `json:"sort_by"`
	Order  string `json:"order"`
}
