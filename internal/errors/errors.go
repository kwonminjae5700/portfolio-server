package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code=%d, message=%s, detail=%s", e.Code, e.Message, e.Detail)
}

func NewAppError(code int, message string, detail string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

var (
	ErrBadRequest          = NewAppError(http.StatusBadRequest, "Bad Request", "")
	ErrUnauthorized        = NewAppError(http.StatusUnauthorized, "Unauthorized", "")
	ErrForbidden           = NewAppError(http.StatusForbidden, "Forbidden", "")
	ErrNotFound            = NewAppError(http.StatusNotFound, "Not Found", "")
	ErrConflict            = NewAppError(http.StatusConflict, "Conflict", "")
	ErrInternalServerError = NewAppError(http.StatusInternalServerError, "Internal Server Error", "")
)

func ErrUserNotFound() *AppError {
	return NewAppError(http.StatusNotFound, "User not found", "The requested user does not exist")
}

func ErrArticleNotFound() *AppError {
	return NewAppError(http.StatusNotFound, "Article not found", "The requested article does not exist")
}

func ErrCommentNotFound() *AppError {
	return NewAppError(http.StatusNotFound, "Comment not found", "The requested comment does not exist")
}

func ErrInvalidCredentials() *AppError {
	return NewAppError(http.StatusUnauthorized, "Invalid credentials", "Email or password is incorrect")
}

func ErrEmailAlreadyExists() *AppError {
	return NewAppError(http.StatusConflict, "Email already exists", "A user with this email already exists")
}

func ErrUsernameAlreadyExists() *AppError {
	return NewAppError(http.StatusConflict, "Username already exists", "A user with this username already exists")
}

func ErrInvalidInput(detail string) *AppError {
	return NewAppError(http.StatusBadRequest, "Invalid input", detail)
}

func ErrPermissionDenied() *AppError {
	return NewAppError(http.StatusForbidden, "Permission denied", "You don't have permission to perform this action")
}
