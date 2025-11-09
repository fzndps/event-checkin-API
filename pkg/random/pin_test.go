package random

import (
	"fmt"
	"testing"
)

func TestGeneratePIN(t *testing.T) {
	// Generate multiple PINs
	pins := make(map[string]bool)

	for i := 0; i < 100; i++ {
		pin, err := GeneratePIN()
		if err != nil {
			t.Fatal("Failed to generate PIN:", err)
		}

		// Check length (must be 4 digits)
		if len(pin) != 4 {
			t.Errorf("PIN length should be 4, got %d: %s", len(pin), pin)
		}

		// Check format (must be numeric)
		for _, char := range pin {
			if char < '0' || char > '9' {
				t.Errorf("PIN should contain only digits, got: %s", pin)
			}
		}

		// Collect PINs to check uniqueness
		pins[pin] = true
	}
	fmt.Println(pins)
	// Check that we got some variety (not all same)
	if len(pins) < 50 {
		t.Errorf("Expected more unique PINs, got only %d out of 100", len(pins))
	}

	t.Logf("✅ Generated %d unique PINs out of 100", len(pins))
}

func TestGeneratePIN_LeadingZeros(t *testing.T) {
	// Generate many PINs to ensure leading zeros are preserved
	for i := 0; i < 1000; i++ {
		pin, err := GeneratePIN()
		if err != nil {
			t.Fatal("Failed to generate PIN:", err)
		}
		fmt.Println(pin)
		if len(pin) != 4 {
			t.Errorf("PIN should always be 4 digits (with leading zeros): %s", pin)
		}
	}

	t.Log("✅ All PINs have correct length with leading zeros")
}

func TestGenerateToken(t *testing.T) {
	// Generate multiple tokens
	tokens := make(map[string]bool)

	for i := 0; i < 100; i++ {
		token, err := GenerateToken()
		if err != nil {
			t.Fatal("Failed to generate token:", err)
		}

		// Check length (must be 32 hex characters)
		if len(token) != 32 {
			t.Errorf("Token length should be 32, got %d: %s", len(token), token)
		}

		// Check format (must be hexadecimal)
		for _, char := range token {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
				t.Errorf("Token should contain only hex characters, got: %s", token)
			}
		}

		// Check uniqueness
		if tokens[token] {
			t.Error("Generated duplicate token!")
		}
		tokens[token] = true
	}

	fmt.Println(tokens)
	// All tokens should be unique
	if len(tokens) != 100 {
		t.Errorf("Expected 100 unique tokens, got %d", len(tokens))
	}

	t.Log("✅ All tokens are unique and properly formatted")
}

func TestGenerateToken_Collision(t *testing.T) {
	// Generate many tokens to test collision resistance
	tokens := make(map[string]bool)
	iterations := 10000

	for i := 0; i < iterations; i++ {
		token, err := GenerateToken()
		if err != nil {
			t.Fatal("Failed to generate token:", err)
		}

		if tokens[token] {
			t.Fatalf("Collision detected after %d iterations!", i)
		}
		tokens[token] = true
	}

	fmt.Println(tokens)

	t.Logf("✅ No collisions in %d tokens", iterations)
}

func BenchmarkGeneratePIN(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GeneratePIN()
	}
}

func BenchmarkGenerateToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateToken()
	}
}
