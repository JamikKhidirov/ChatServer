package verservice

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	authdomain "ChatServerGolang/backend/internal/domain/auth"
	verificationdomain "ChatServerGolang/backend/internal/domain/verification"
	"ChatServerGolang/backend/internal/email"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/internal/sms"

	"github.com/google/uuid"
)

type verificationService struct {
	verRepo     repository.VerificationRepository
	userRepo    repository.UserRepository
	emailSender *email.Sender
	smsSender   *sms.Sender
}

func NewVerificationService(verRepo repository.VerificationRepository, userRepo repository.UserRepository, emailSender *email.Sender, smsSender *sms.Sender) service.VerificationService {
	return &verificationService{verRepo: verRepo, userRepo: userRepo, emailSender: emailSender, smsSender: smsSender}
}

func generateCode() string {
	code, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", code.Int64()+100000)
}

func (s *verificationService) SendEmailVerification(userID, email string) error {
	code := generateCode()
	ver := &verificationdomain.EmailVerification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute).Format(time.RFC3339),
	}
	if err := s.verRepo.CreateEmail(ver); err != nil {
		return err
	}
	body := fmt.Sprintf("Your email verification code is: %s\n\nThis code expires in 10 minutes.", code)
	if err := s.emailSender.Send(email, "Verify Your Email", body); err != nil {
		fmt.Printf("[EMAIL VERIFICATION FAILED] User %s code: %s (to: %s) - %v\n", userID, code, email, err)
	}
	fmt.Printf("[EMAIL VERIFICATION] User %s code: %s (to: %s)\n", userID, code, email)
	return nil
}

func (s *verificationService) VerifyEmail(userID, code string) error {
	ver, err := s.verRepo.FindEmailByUserID(userID)
	if err != nil {
		return fmt.Errorf("verification not found")
	}
	if ver.Verified == 1 {
		return fmt.Errorf("already verified")
	}
	expires, _ := time.Parse(time.RFC3339, ver.ExpiresAt)
	if time.Now().After(expires) {
		return fmt.Errorf("code expired")
	}
	if ver.Code != code {
		return fmt.Errorf("invalid code")
	}
	return s.verRepo.VerifyEmail(ver.ID)
}

func (s *verificationService) SendPhoneVerification(userID, phone string) error {
	code := generateCode()
	ver := &verificationdomain.PhoneVerification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Phone:     phone,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute).Format(time.RFC3339),
	}
	if err := s.verRepo.CreatePhone(ver); err != nil {
		return err
	}
	body := fmt.Sprintf("Your phone verification code is: %s", code)
	if err := s.smsSender.Send(phone, body); err != nil {
		fmt.Printf("[SMS VERIFICATION FAILED] User %s code: %s (to: %s) - %v\n", userID, code, phone, err)
	}
	fmt.Printf("[SMS VERIFICATION] User %s code: %s (to: %s)\n", userID, code, phone)
	return nil
}

func (s *verificationService) VerifyPhone(userID, code string) error {
	ver, err := s.verRepo.FindPhoneByUserID(userID)
	if err != nil {
		return fmt.Errorf("verification not found")
	}
	if ver.Verified == 1 {
		return fmt.Errorf("already verified")
	}
	expires, _ := time.Parse(time.RFC3339, ver.ExpiresAt)
	if time.Now().After(expires) {
		return fmt.Errorf("code expired")
	}
	if ver.Code != code {
		return fmt.Errorf("invalid code")
	}
	return s.verRepo.VerifyPhone(ver.ID)
}

func (s *verificationService) IsEmailVerified(userID string) (bool, error) {
	ver, err := s.verRepo.FindEmailByUserID(userID)
	if err != nil {
		return false, nil
	}
	return ver.Verified == 1, nil
}

func (s *verificationService) IsPhoneVerified(userID string) (bool, error) {
	ver, err := s.verRepo.FindPhoneByUserID(userID)
	if err != nil {
		return false, nil
	}
	return ver.Verified == 1, nil
}

func (s *verificationService) LoginSendEmailCode(email string) (string, error) {
	_, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("email not registered")
	}
	code := generateCode()
	loginCode := &authdomain.EmailLoginCode{
		ID:        uuid.New().String(),
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute).Format(time.RFC3339),
	}
	if err := s.verRepo.CreateEmailLoginCode(loginCode); err != nil {
		return "", err
	}
	body := fmt.Sprintf("Your login code is: %s\n\nThis code expires in 10 minutes.", code)
	if err := s.emailSender.Send(email, "Your Login Code", body); err != nil {
		fmt.Printf("[EMAIL LOGIN CODE FAILED] %s - %v\n", email, err)
	}
	fmt.Printf("[EMAIL LOGIN CODE] %s code: %s\n", email, code)
	return code, nil
}

func (s *verificationService) LoginVerifyEmailCode(email, code string) (string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("email not registered")
	}
	loginCode, err := s.verRepo.FindEmailLoginCode(email)
	if err != nil {
		return "", errors.New("no verification code found")
	}
	if loginCode.Verified == 1 {
		return "", errors.New("code already used")
	}
	expires, _ := time.Parse(time.RFC3339, loginCode.ExpiresAt)
	if time.Now().After(expires) {
		return "", errors.New("code expired")
	}
	if loginCode.Code != code {
		return "", errors.New("invalid code")
	}
	if err := s.verRepo.VerifyEmailLoginCode(loginCode.ID); err != nil {
		return "", err
	}
	return user.ID, nil
}

func (s *verificationService) LoginSendPhoneCode(phone string) (string, error) {
	_, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return "", errors.New("phone not registered")
	}
	code := generateCode()
	loginCode := &authdomain.PhoneLoginCode{
		ID:        uuid.New().String(),
		Phone:     phone,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute).Format(time.RFC3339),
	}
	if err := s.verRepo.CreatePhoneLoginCode(loginCode); err != nil {
		return "", err
	}
	body := fmt.Sprintf("Your login code is: %s", code)
	if err := s.smsSender.Send(phone, body); err != nil {
		fmt.Printf("[SMS LOGIN CODE FAILED] %s - %v\n", phone, err)
	}
	fmt.Printf("[PHONE LOGIN CODE] %s code: %s\n", phone, code)
	return code, nil
}

func (s *verificationService) LoginVerifyPhoneCode(phone, code string) (string, error) {
	user, err := s.userRepo.FindByPhone(phone)
	if err != nil {
		return "", errors.New("phone not registered")
	}
	loginCode, err := s.verRepo.FindPhoneLoginCode(phone)
	if err != nil {
		return "", errors.New("no verification code found")
	}
	if loginCode.Verified == 1 {
		return "", errors.New("code already used")
	}
	expires, _ := time.Parse(time.RFC3339, loginCode.ExpiresAt)
	if time.Now().After(expires) {
		return "", errors.New("code expired")
	}
	if loginCode.Code != code {
		return "", errors.New("invalid code")
	}
	if err := s.verRepo.VerifyPhoneLoginCode(loginCode.ID); err != nil {
		return "", err
	}
	return user.ID, nil
}


