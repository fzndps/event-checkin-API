package jwt

import (
	"fmt"
	"testing"
	"time"
)

var secretKey = "test-secret-key-min-32-character"

var tokenDuration = 1 * time.Hour

func TestGenerateToken(t *testing.T) {
	manager := NewJWTManager(secretKey, tokenDuration)

	token, err := manager.GenerateToken(1, "test@example.com")
	if err != nil {
		t.Fatal("Failed to generate token:", err)
	}

	if token == "" {
		t.Fatal("Token should not be empt")
	}

	t.Log("Token generated:", token)
}

func TestTokenValidate(t *testing.T) {
	manager := NewJWTManager(secretKey, tokenDuration)

	token, _ := manager.GenerateToken(123, "test@example.com")

	claims, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatal("Failed to validate token:", err)
	}

	if claims.OrganizerID != 123 {
		t.Fatalf("Expected organizer_id 123, got %d", claims.OrganizerID)
	}

	if claims.Email != "test@example.com" {
		t.Fatalf("Expected email test@example.com, got %s", claims.Email)
	}

	fmt.Println(claims.RegisteredClaims.ExpiresAt)
	fmt.Println(token)

	t.Log("Token validate successfully")
}

func TestExpiredToken(t *testing.T) {
	manager := NewJWTManager(secretKey, 1*time.Nanosecond)

	token, _ := manager.GenerateToken(1, "test@example.com")

	time.Sleep(10 * time.Millisecond)

	_, err := manager.ValidateToken(token)
	if err == nil {
		t.Fatal("Should return error for expired token")
	}

	t.Log("Expired token rejected correctly")
}
