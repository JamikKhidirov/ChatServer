package domain

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
