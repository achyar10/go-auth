package auth

type RegisterDTO struct {
	Username string `json:"username" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required,min=6"`
	Fullname string `json:"fullname"`
	Role     string `json:"role" validate:"oneof=admin user"`
}

type LoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
