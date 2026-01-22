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
	ErrBadRequest          = NewAppError(http.StatusBadRequest, "잘못된 요청입니다", "")
	ErrUnauthorized        = NewAppError(http.StatusUnauthorized, "인증이 필요합니다", "")
	ErrForbidden           = NewAppError(http.StatusForbidden, "접근 권한이 없습니다", "")
	ErrNotFound            = NewAppError(http.StatusNotFound, "요청한 리소스를 찾을 수 없습니다", "")
	ErrConflict            = NewAppError(http.StatusConflict, "중복된 요청입니다", "")
	ErrInternalServerError = NewAppError(http.StatusInternalServerError, "서버 오류가 발생했습니다", "")
)

func ErrUserNotFound() *AppError {
	return NewAppError(http.StatusNotFound, "사용자를 찾을 수 없습니다", "요청한 사용자가 존재하지 않습니다")
}

func ErrArticleNotFound() *AppError {
	return NewAppError(http.StatusNotFound, "게시글을 찾을 수 없습니다", "요청한 게시글이 존재하지 않습니다")
}

func ErrCommentNotFound() *AppError {
	return NewAppError(http.StatusNotFound, "댓글을 찾을 수 없습니다", "요청한 댓글이 존재하지 않습니다")
}

func ErrInvalidCredentials() *AppError {
	return NewAppError(http.StatusUnauthorized, "로그인 정보가 올바르지 않습니다", "이메일 또는 비밀번호가 일치하지 않습니다")
}

func ErrEmailAlreadyExists() *AppError {
	return NewAppError(http.StatusConflict, "이미 사용 중인 이메일입니다", "해당 이메일로 가입된 계정이 이미 존재합니다")
}

func ErrUsernameAlreadyExists() *AppError {
	return NewAppError(http.StatusConflict, "이미 사용 중인 사용자명입니다", "해당 사용자명은 이미 다른 사용자가 사용하고 있습니다")
}

func ErrInvalidInput(detail string) *AppError {
	return NewAppError(http.StatusBadRequest, "입력 값이 올바르지 않습니다", detail)
}

func ErrPermissionDenied() *AppError {
	return NewAppError(http.StatusForbidden, "권한이 없습니다", "이 작업을 수행할 권한이 없습니다")
}
