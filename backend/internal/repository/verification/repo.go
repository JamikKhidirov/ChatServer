package verrepo

import (
	"database/sql"
	"time"

	authdomain "ChatServerGolang/backend/internal/domain/auth"
	verificationdomain "ChatServerGolang/backend/internal/domain/verification"
	"ChatServerGolang/backend/internal/repository"
)

type verificationRepository struct {
	db *sql.DB
}

func NewVerificationRepository(db *sql.DB) repository.VerificationRepository {
	return &verificationRepository{db: db}
}

func (r *verificationRepository) CreateEmail(ver *verificationdomain.EmailVerification) error {
	_, err := r.db.Exec(`INSERT INTO email_verifications (id, user_id, email, code, expires_at, verified, created_at) VALUES (?, ?, ?, ?, ?, 0, ?)`,
		ver.ID, ver.UserID, ver.Email, ver.Code, ver.ExpiresAt, time.Now().Format(time.RFC3339))
	return err
}

func (r *verificationRepository) FindEmailByUserID(userID string) (*verificationdomain.EmailVerification, error) {
	row := r.db.QueryRow(`SELECT id, user_id, email, code, expires_at, verified, created_at FROM email_verifications WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`, userID)
	var ver verificationdomain.EmailVerification
	err := row.Scan(&ver.ID, &ver.UserID, &ver.Email, &ver.Code, &ver.ExpiresAt, &ver.Verified, &ver.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &ver, nil
}

func (r *verificationRepository) VerifyEmail(id string) error {
	_, err := r.db.Exec(`UPDATE email_verifications SET verified = 1 WHERE id = ?`, id)
	return err
}

func (r *verificationRepository) CreatePhone(ver *verificationdomain.PhoneVerification) error {
	_, err := r.db.Exec(`INSERT INTO phone_verifications (id, user_id, phone, code, expires_at, verified, created_at) VALUES (?, ?, ?, ?, ?, 0, ?)`,
		ver.ID, ver.UserID, ver.Phone, ver.Code, ver.ExpiresAt, time.Now().Format(time.RFC3339))
	return err
}

func (r *verificationRepository) FindPhoneByUserID(userID string) (*verificationdomain.PhoneVerification, error) {
	row := r.db.QueryRow(`SELECT id, user_id, phone, code, expires_at, verified, created_at FROM phone_verifications WHERE user_id = ? ORDER BY created_at DESC LIMIT 1`, userID)
	var ver verificationdomain.PhoneVerification
	err := row.Scan(&ver.ID, &ver.UserID, &ver.Phone, &ver.Code, &ver.ExpiresAt, &ver.Verified, &ver.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &ver, nil
}

func (r *verificationRepository) VerifyPhone(id string) error {
	_, err := r.db.Exec(`UPDATE phone_verifications SET verified = 1 WHERE id = ?`, id)
	return err
}

func (r *verificationRepository) CreateEmailLoginCode(code *authdomain.EmailLoginCode) error {
	_, err := r.db.Exec(`INSERT INTO email_login_codes (id, email, code, expires_at, verified, created_at) VALUES (?, ?, ?, ?, 0, ?)`,
		code.ID, code.Email, code.Code, code.ExpiresAt, time.Now().Format(time.RFC3339))
	return err
}

func (r *verificationRepository) FindEmailLoginCode(email string) (*authdomain.EmailLoginCode, error) {
	row := r.db.QueryRow(`SELECT id, email, code, expires_at, verified, created_at FROM email_login_codes WHERE email = ? ORDER BY created_at DESC LIMIT 1`, email)
	var code authdomain.EmailLoginCode
	err := row.Scan(&code.ID, &code.Email, &code.Code, &code.ExpiresAt, &code.Verified, &code.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &code, nil
}

func (r *verificationRepository) VerifyEmailLoginCode(id string) error {
	_, err := r.db.Exec(`UPDATE email_login_codes SET verified = 1 WHERE id = ?`, id)
	return err
}

func (r *verificationRepository) CreatePhoneLoginCode(code *authdomain.PhoneLoginCode) error {
	_, err := r.db.Exec(`INSERT INTO phone_login_codes (id, phone, code, expires_at, verified, created_at) VALUES (?, ?, ?, ?, 0, ?)`,
		code.ID, code.Phone, code.Code, code.ExpiresAt, time.Now().Format(time.RFC3339))
	return err
}

func (r *verificationRepository) FindPhoneLoginCode(phone string) (*authdomain.PhoneLoginCode, error) {
	row := r.db.QueryRow(`SELECT id, phone, code, expires_at, verified, created_at FROM phone_login_codes WHERE phone = ? ORDER BY created_at DESC LIMIT 1`, phone)
	var code authdomain.PhoneLoginCode
	err := row.Scan(&code.ID, &code.Phone, &code.Code, &code.ExpiresAt, &code.Verified, &code.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &code, nil
}

func (r *verificationRepository) VerifyPhoneLoginCode(id string) error {
	_, err := r.db.Exec(`UPDATE phone_login_codes SET verified = 1 WHERE id = ?`, id)
	return err
}



