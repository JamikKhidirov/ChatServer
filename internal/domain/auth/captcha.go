package authdomain

type CaptchaToken struct {
	Token     string `json:"token"`
	Solution  string `json:"-"`
	ExpiresAt string `json:"expires_at"`
	Used      int    `json:"used"`
}

type CaptchaResponse struct {
	Token     string `json:"token"`
	Question  string `json:"question"`
	ImageData string `json:"image_data,omitempty"`
}

type CaptchaVerifyRequest struct {
	Token    string `json:"token" binding:"required"`
	Solution string `json:"solution" binding:"required"`
}
