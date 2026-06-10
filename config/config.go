package config

import (
	"os"
	"strconv"
)

type Config struct {
	ServerPort   string
	DatabasePath  string
	JWTSecret    string
	JWTTTL       int
	AllowOrigins  string
	PushEnabled  bool
	FirebaseCredentials string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPass     string
	SMTPFrom     string
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioPhone      string
}

func Load() *Config {
	return &Config{
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DatabasePath: getEnv("DATABASE_PATH", "file:chat.db?cache=shared&mode=rwc"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		JWTTTL:      getEnvInt("JWT_TTL", 86400),
		AllowOrigins: getEnv("ALLOW_ORIGINS", "*"),
		PushEnabled: getEnvBool("PUSH_ENABLED", false),
		FirebaseCredentials: getEnv("FIREBASE_CREDENTIALS", ""),
		SMTPHost:    getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:    getEnv("SMTP_PORT", "587"),
		SMTPUser:    getEnv("SMTP_USER", ""),
		SMTPPass:    getEnv("SMTP_PASS", ""),
		SMTPFrom:    getEnv("SMTP_FROM", "noreply@chatmessenger.local"),
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioPhone:      getEnv("TWILIO_PHONE", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
