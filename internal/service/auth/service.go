package authservice

import (
	"errors"
	"time"

	"ChatServerGolang/internal/config"
	"ChatServerGolang/internal/domain/auth"
	userdomain "ChatServerGolang/internal/domain/user"
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/service"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) service.AuthService {
	return &authService{userRepo: userRepo, cfg: cfg}
}

func (s *authService) Register(req *authdomain.RegisterRequest) (*authdomain.AuthResponse, error) {
	return s.registerUser(req.Username, req.Email, req.Password, req.DisplayName, false)
}

func (s *authService) RegisterAdmin(req *authdomain.AdminRegisterRequest) (*authdomain.AuthResponse, error) {
	if req.AdminSecret != s.cfg.AdminSecret {
		return nil, errors.New("invalid admin secret")
	}
	return s.registerUser(req.Username, req.Email, req.Password, req.DisplayName, true)
}

func (s *authService) registerUser(username, email, password, displayName string, isAdmin bool) (*authdomain.AuthResponse, error) {
	existing, err := s.userRepo.FindByEmail(email)
	if err == nil && existing != nil {
		return nil, errors.New("email already registered")
	}

	existing, err = s.userRepo.FindByUsername(username)
	if err == nil && existing != nil {
		return nil, errors.New("username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &userdomain.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		DisplayName:  displayName,
		Status:       "Available",
		LastSeen:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
		IsAdmin:      isAdmin,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &authdomain.AuthResponse{Token: token, User: user}, nil
}

func (s *authService) Login(req *authdomain.LoginRequest) (*authdomain.AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &authdomain.AuthResponse{Token: token, User: user}, nil
}

func (s *authService) RefreshToken(userID string) (string, error) {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}
	return s.generateToken(userID)
}

func (s *authService) ChangePassword(userID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, string(hash))
}

func (s *authService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	return userID, nil
}

func (s *authService) generateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(s.cfg.JWTTTL) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}


