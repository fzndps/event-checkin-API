// Package middleware berfungsi menjalankan middleware sebelum handler
package middleware

import (
	"strings"

	"github.com/fzndps/eventcheck/pkg/jwt"
	"github.com/fzndps/eventcheck/pkg/validator"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtManager *jwt.JWTManager
}

func NewAuthMiddleware(jwtManager *jwt.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	// Return function akan dijalankan saat ada request
	return func(c *gin.Context) {
		// authHeader ini bisasanya berbentuk : Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			validator.UnauthorizedResponse(c, "Authorization header required")
			c.Abort()
			return
		}

		// Split berfungsi untuk mengubah "Bearer <token>" menjadi ["Bearer", "<token>"]
		parts := strings.Split(authHeader, " ")

		// Validasi format harus ada 2 parts dan parts pertama harus berisi "Bearer"
		if len(parts) != 2 || parts[0] != "Bearer" {
			validator.UnauthorizedResponse(c, "Invalid authorization format.")
			c.Abort()
			return
		}

		// Parts ke 2 adalah token
		tokenString := parts[1]

		// ValidateJWT dari validator akan mengecek token masih valid dan belum expired
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			validator.UnauthorizedResponse(c, err.Error())
			c.Abort()
			return
		}

		// c.Set untuk menyimpan data ke context dan bisa di ambil
		// dihandler untuk mengetahui siapa yang login dengan c.Get
		c.Set("organizer_id", claims.OrganizerID)
		c.Set("organizer_email", claims.Email)

		// Lanjut ke middleware/handler berikutnya
		c.Next()
	}
}

// GetUserID untuk mengambil userID
func GetOrganizerID(c *gin.Context) (int64, bool) {
	// Ambil data dari context
	organizerID, exists := c.Get("organizer_id")
	if !exists {
		return 0, false
	}

	// Convert nilai interface/any ke int
	id, ok := organizerID.(int64)
	if !ok {
		return 0, false
	}

	return id, true
}

func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("organizer_email")
	if !exists {
		return "", false
	}

	emailSTR, ok := email.(string)
	if !ok {
		return "", false
	}

	return emailSTR, true
}
