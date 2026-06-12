package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Total  int `json:"total,omitempty"`
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

type PaginatedData struct {
	Items interface{} `json:"items"`
	Total int         `json:"total"`
}

const (
	ErrBadRequest          = "BAD_REQUEST"
	ErrUnauthorized        = "UNAUTHORIZED"
	ErrForbidden           = "FORBIDDEN"
	ErrNotFound            = "NOT_FOUND"
	ErrInternal            = "INTERNAL_ERROR"
	ErrValidation          = "VALIDATION_ERROR"
	ErrDuplicate           = "DUPLICATE"
	ErrRateLimit           = "RATE_LIMIT"
	ErrBlocked             = "BLOCKED"
	ErrAccessDenied        = "ACCESS_DENIED"
	ErrUserNotFound        = "USER_NOT_FOUND"
	ErrChatNotFound        = "CHAT_NOT_FOUND"
	ErrMessageNotFound     = "MESSAGE_NOT_FOUND"
	ErrInvalidToken        = "INVALID_TOKEN"
	ErrTokenExpired        = "TOKEN_EXPIRED"
	ErrUsernameTaken       = "USERNAME_TAKEN"
	ErrEmailTaken          = "EMAIL_TAKEN"
	ErrWeakPassword        = "WEAK_PASSWORD"
)

func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, APIResponse{Success: true, Data: data})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{Success: false, Error: message})
}

func ErrorWithCode(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{Success: false, Error: message, Code: code})
}

func BadRequest(c *gin.Context, message string) {
	ErrorWithCode(c, http.StatusBadRequest, ErrBadRequest, message)
}

func ValidationError(c *gin.Context, message string) {
	ErrorWithCode(c, http.StatusBadRequest, ErrValidation, message)
}

func Unauthorized(c *gin.Context, message string) {
	ErrorWithCode(c, http.StatusUnauthorized, ErrUnauthorized, message)
}

func NotFound(c *gin.Context, message string) {
	ErrorWithCode(c, http.StatusNotFound, ErrNotFound, message)
}

func InternalError(c *gin.Context, message string) {
	ErrorWithCode(c, http.StatusInternalServerError, ErrInternal, message)
}

func Forbidden(c *gin.Context, message string) {
	ErrorWithCode(c, http.StatusForbidden, ErrForbidden, message)
}

func Conflict(c *gin.Context, code, message string) {
	ErrorWithCode(c, http.StatusConflict, code, message)
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Error   string `json:"error" example:"error description"`
	Code    string `json:"code,omitempty" example:"ERROR_CODE"`
}

type MessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Data    gin.H  `json:"data,omitempty"`
}

func Paginated(c *gin.Context, status int, items interface{}, total, limit, offset int) {
	c.JSON(status, APIResponse{
		Success: true,
		Data:    items,
		Meta: &Meta{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	})
}
