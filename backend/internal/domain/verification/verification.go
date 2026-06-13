package verificationdomain

type EmailVerification struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Code      string `json:"-"`
	ExpiresAt string `json:"expires_at"`
	Verified  int    `json:"verified"`
	CreatedAt string `json:"created_at"`
}

type PhoneVerification struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Phone     string `json:"phone"`
	Code      string `json:"-"`
	ExpiresAt string `json:"expires_at"`
	Verified  int    `json:"verified"`
	CreatedAt string `json:"created_at"`
}

type SendEmailVerificationRequest struct {
	Email string `json:"email" binding:"required"`
}

type VerifyEmailRequest struct {
	Code string `json:"code" binding:"required"`
}

type SendPhoneVerificationRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type VerifyPhoneRequest struct {
	Code string `json:"code" binding:"required"`
}
