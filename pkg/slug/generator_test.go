package slug

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple text",
			input:    "Tech Conference 2024",
			expected: "tech-conference-2024",
		},
		{
			name:     "Text with special characters",
			input:    "Hello World!!!",
			expected: "hello-world",
		},
		{
			name:     "Text with multiple spaces",
			input:    "   Spasi   Banyak   ",
			expected: "spasi-banyak",
		},
		{
			name:     "Text with dashes",
			input:    "Already-Has-Dashes",
			expected: "already-has-dashes",
		},
		{
			name:     "Text with unicode characters",
			input:    "Café Résumé",
			expected: "caf-rsum",
		},
		{
			name:     "Mixed case",
			input:    "MiXeD CaSe TeXt",
			expected: "mixed-case-text",
		},
		{
			name:     "Numbers and text",
			input:    "Event 123 456",
			expected: "event-123-456",
		},
		{
			name:     "Multiple consecutive dashes",
			input:    "Text---With---Dashes",
			expected: "text-with-dashes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Generate(tt.input)
			fmt.Println("result: ", result)
			if result != tt.expected {
				t.Errorf("Generate(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerate_MaxLength(t *testing.T) {
	// Create very long input
	longInput := strings.Repeat("a", 150)

	result := Generate(longInput)
	fmt.Println("result: ", result)
	if len(result) > 100 {
		t.Errorf("Generate() returned slug longer than 100 chars: %d", len(result))
	}
}

func TestGenerateUnique(t *testing.T) {
	input := "Tech Conference"

	// Generate 2 unique slugs
	slug1 := GenerateUnique(input)
	time.Sleep(1 * time.Second)
	slug2 := GenerateUnique(input)

	// Should be different (contains timestamp)
	if slug1 == slug2 {
		t.Error("GenerateUnique() should return different slugs")
	}

	fmt.Println("Slug 1: ", slug1)
	fmt.Println("Slug 2: ", slug2)
	// Both should start with base slug
	baseSlug := Generate(input)
	if !strings.HasPrefix(slug1, baseSlug) {
		t.Errorf("Unique slug should start with base slug: %s", slug1)
	}
	if !strings.HasPrefix(slug2, baseSlug) {
		t.Errorf("Unique slug should start with base slug: %s", slug2)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		expected bool
	}{
		{
			name:     "Valid slug",
			slug:     "tech-conference-2024",
			expected: true,
		},
		{
			name:     "Valid slug with numbers",
			slug:     "event-123",
			expected: true,
		},
		{
			name:     "Empty slug",
			slug:     "",
			expected: false,
		},
		{
			name:     "Slug with uppercase",
			slug:     "Tech-Conference",
			expected: false,
		},
		{
			name:     "Slug with spaces",
			slug:     "tech conference",
			expected: false,
		},
		{
			name:     "Slug with special chars",
			slug:     "tech-conference!",
			expected: false,
		},
		{
			name:     "Slug too long",
			slug:     strings.Repeat("a", 101),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.slug)
			if result != tt.expected {
				t.Errorf("Validate(%q) = %v, want %v", tt.slug, result, tt.expected)
			}
		})
	}
}

func BenchmarkGenerate(b *testing.B) {
	input := "Tech Conference 2024"

	for i := 0; i < b.N; i++ {
		Generate(input)
	}
}

func BenchmarkGenerateUnique(b *testing.B) {
	input := "Tech Conference 2024"

	for i := 0; i < b.N; i++ {
		GenerateUnique(input)
	}
}
