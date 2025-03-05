package auth

type LoginResponse struct {
	Id       int64  `json:"user_id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Role     string `json:"role"`
	Token    string `json:"access_token"`
}
