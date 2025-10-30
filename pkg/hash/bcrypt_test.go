package hash

import "testing"

// Test hash password
func TestHashPassword(t *testing.T) {
	password := "superSecretPassword12345"

	hashedPass, err := HashPassword(password)
	if err != nil {
		t.Fatal("Failed to hash password:", err)
	}

	if hashedPass == "" {
		t.Fatal("Hash should be not empty")
	}

	if hashedPass == password {
		t.Fatal("Hash should not equal plaintext password")
	}

	t.Log("Password hashed successfully")
}

// Test check password
func TestCheckPassword(t *testing.T) {
	password := "superSecretPassword12345"
	wrongPassword := "anjay123"

	hash, _ := HashPassword(password)

	if !CheckPassword(password, hash) {
		t.Fatal("Should return true for correct password")
	}

	if CheckPassword(wrongPassword, hash) {
		t.Fatal("Should return false for wrong password")
	}

	t.Log("Password checking working correctly")
}
