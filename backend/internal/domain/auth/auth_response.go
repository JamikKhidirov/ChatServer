package authdomain

import userdomain "ChatServerGolang/backend/internal/domain/user"

type AuthResponse struct {
	Token string          `json:"token"`
	User  *userdomain.User `json:"user"`
}
