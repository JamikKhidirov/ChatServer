package domain

type EmailLoginCode struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Code      string `json:"-"`
	ExpiresAt string `json:"expires_at"`
	Verified  int    `json:"verified"`
	CreatedAt string `json:"created_at"`
}

type PhoneLoginCode struct {
	ID        string `json:"id"`
	Phone     string `json:"phone"`
	Code      string `json:"-"`
	ExpiresAt string `json:"expires_at"`
	Verified  int    `json:"verified"`
	CreatedAt string `json:"created_at"`
}

type LoginByEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

type LoginByEmailVerifyRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type LoginByPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type LoginByPhoneVerifyRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}
