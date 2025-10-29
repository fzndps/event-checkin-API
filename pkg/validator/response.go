// Package validator response untuk membantu menentukan response dari controller/handler
package validator

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// SuccessResponse mengirim response sukses dengan HTTP 200
func SuccessResponse(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// CreatedResponse mengirim response created dengan HTTP 201
func CreatedResponse(c *gin.Context, message string, data any) {
	c.JSON(http.StatusCreated, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengirim response error dengan HTTP custom
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Status:  false,
		Message: "Error",
		Error:   message,
	})
}

// BadRequestResponse mengirim response bad request dengan HTTP 400
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// UnauthorizedResponse mengirim response Unauthorized dengan HTTP 401
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message)
}

// ForbiddenResponse mengirim response Unauthorized dengan HTTP 403
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message)
}

// NotFoundResponse mengirim response Unauthorized dengan HTTP 404
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

// InternalServerErrorResponse mengirim response Unauthorized dengan HTTP 500
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}

func ValidateErrorResponse(c *gin.Context, err any) {
	c.JSON(http.StatusBadRequest, Response{
		Status:  false,
		Message: "Validation error",
		Data:    err,
	})
}
