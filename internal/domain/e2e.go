package domain

type E2EKey struct {
	UserID             string `json:"user_id"`
	PublicKey          string `json:"public_key"`
	PrivateKeyEncrypted string `json:"-"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type E2EAdminKey struct {
	ID              int    `json:"id"`
	AdminPublicKey  string `json:"admin_public_key"`
	AdminPrivateKey string `json:"-"`
	CreatedAt       string `json:"created_at"`
}

type E2ERegisterRequest struct {
	PublicKey          string `json:"public_key" binding:"required"`
	PrivateKeyEncrypted string `json:"private_key_encrypted" binding:"required"`
}

type E2EEncryptRequest struct {
	Message  string `json:"message" binding:"required"`
	ChatID   string `json:"chat_id" binding:"required"`
}

type E2EDecryptRequest struct {
	Ciphertext string `json:"ciphertext" binding:"required"`
	ChatID     string `json:"chat_id" binding:"required"`
}

type E2EEncryptedMessage struct {
	ChatID          string `json:"chat_id"`
	SenderID        string `json:"sender_id"`
	Ciphertext      string `json:"ciphertext"`
	EncryptedKey    string `json:"encrypted_key"`
	Nonce           string `json:"nonce"`
	CreatedAt       string `json:"created_at"`
}
