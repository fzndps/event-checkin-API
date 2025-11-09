package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	OrganizerID int    `json:"organizer_id"`
	Email       string `json:"email"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey string
}

func NewJWTManager(secretKey string) *JWTManager {
	return &JWTManager{
		secretKey: secretKey,
	}
}

func (m *JWTManager) GenerateToken(organizerID int, email string, expiryHours int) (string, error) {

	if len(m.secretKey) == 0 {
		return "", errors.New("JWT secret not initialize")
	}

	claims := &JWTClaims{
		OrganizerID: organizerID,
		Email:       email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiryHours))),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {

	if len(m.secretKey) == 0 {
		return nil, errors.New("JWT secret not initialize")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (m *JWTManager) RefreshToken(oldTokenString string, expiryHours int) (string, error) {
	claims, err := m.ValidateToken(oldTokenString)
	if err != nil {
		return "", err
	}

	return m.GenerateToken(claims.OrganizerID, claims.Email, expiryHours)
}
