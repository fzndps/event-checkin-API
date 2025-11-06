package csv

import (
	"strings"
	"testing"

	"github.com/fzndps/eventcheck/internal/domain"
)

func TestParseParticipants_ValidCSV(t *testing.T) {
	csvData := `name,email,phone
John Doe,john@example.com,08123456789
Jane Smith,jane@example.com,08987654321
Bob Johnson,bob@example.com,08111222333`

	reader := strings.NewReader(csvData)
	participants, errors, err := ParseParticipants(reader)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if len(participants) != 3 {
		t.Errorf("Expected 3 participants, got %d", len(participants))
	}

	if len(errors) != 0 {
		t.Errorf("Expected no errors, got %d: %v", len(errors), errors)
	}

	// Verify first participant
	if participants[0].Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", participants[0].Name)
	}
	if participants[0].Email != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%s'", participants[0].Email)
	}
	if participants[0].Phone != "08123456789" {
		t.Errorf("Expected phone '08123456789', got '%s'", participants[0].Phone)
	}
}

func TestParseParticipants_FlexibleColumnOrder(t *testing.T) {
	// Columns in different order
	csvData := `email,phone,name
john@example.com,08123456789,John Doe
jane@example.com,08987654321,Jane Smith`

	reader := strings.NewReader(csvData)
	participants, _, err := ParseParticipants(reader)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if len(participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(participants))
	}

	// Verify data mapped correctly despite different order
	if participants[0].Name != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%s'", participants[0].Name)
	}
}

func TestParseParticipants_MissingColumn(t *testing.T) {
	// Missing 'phone' column
	csvData := `name,email
John Doe,john@example.com`

	reader := strings.NewReader(csvData)
	_, _, err := ParseParticipants(reader)

	if err == nil {
		t.Error("Expected error for missing column, got nil")
	}

	if !strings.Contains(err.Error(), "phone") {
		t.Errorf("Error should mention missing 'phone' column: %v", err)
	}
}

func TestParseParticipants_EmptyFields(t *testing.T) {
	csvData := `name,email,phone
,john@example.com,08123456789
Jane Smith,,08987654321
Bob Johnson,bob@example.com,`

	reader := strings.NewReader(csvData)
	participants, errors, err := ParseParticipants(reader)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	// All rows should be skipped due to empty fields
	if len(participants) != 0 {
		t.Errorf("Expected 0 valid participants, got %d", len(participants))
	}

	// Should have 3 errors
	if len(errors) != 3 {
		t.Errorf("Expected 3 errors, got %d", len(errors))
	}

	t.Logf("Errors: %v", errors)
}

func TestParseParticipants_InvalidEmail(t *testing.T) {
	csvData := `name,email,phone
John Doe,notanemail,08123456789
Jane Smith,jane@example.com,08987654321`

	reader := strings.NewReader(csvData)
	participants, errors, err := ParseParticipants(reader)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	// Only 1 valid participant (Jane)
	if len(participants) != 1 {
		t.Errorf("Expected 1 valid participant, got %d", len(participants))
	}

	// Should have 1 error (invalid email)
	if len(errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(errors))
	}

	if participants[0].Name != "Jane Smith" {
		t.Error("Valid participant should be Jane Smith")
	}
}

func TestParseParticipants_EmptyCSV(t *testing.T) {
	csvData := `name,email,phone`

	reader := strings.NewReader(csvData)
	participants, errors, err := ParseParticipants(reader)

	if err != domain.ErrEmptyCSV {
		t.Errorf("Expected ErrEmptyCSV, got %v", err)
	}

	if participants != nil {
		t.Error("Expected nil participants for empty CSV")
	}

	if len(errors) != 0 {
		t.Error("Expected no errors for empty CSV")
	}
}

func TestParseParticipants_WhitespaceHandling(t *testing.T) {
	csvData := `name,email,phone
  John Doe  ,  john@example.com  ,  08123456789  
Jane Smith,jane@example.com,08987654321`

	reader := strings.NewReader(csvData)
	participants, _, err := ParseParticipants(reader)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if len(participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(participants))
	}

	// Verify whitespace trimmed
	if participants[0].Name != "John Doe" {
		t.Errorf("Expected trimmed name 'John Doe', got '%s'", participants[0].Name)
	}
	if participants[0].Email != "john@example.com" {
		t.Errorf("Expected trimmed email, got '%s'", participants[0].Email)
	}
}

func TestParseParticipants_CaseInsensitiveHeaders(t *testing.T) {
	// Headers with different cases
	csvData := `NAME,EMAIL,PHONE
John Doe,john@example.com,08123456789`

	reader := strings.NewReader(csvData)
	participants, _, err := ParseParticipants(reader)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if len(participants) != 1 {
		t.Errorf("Expected 1 participant, got %d", len(participants))
	}
}

func TestParseParticipants_LargeFile(t *testing.T) {
	// Create CSV with 1000 rows
	var sb strings.Builder
	sb.WriteString("name,email,phone\n")

	for i := 0; i < 1000; i++ {
		sb.WriteString("User")
		sb.WriteString(strings.Repeat(string(rune('0'+i%10)), 1))
		sb.WriteString(",user")
		sb.WriteString(strings.Repeat(string(rune('0'+i%10)), 1))
		sb.WriteString("@example.com,0812345678")
		sb.WriteString(strings.Repeat(string(rune('0'+i%10)), 1))
		sb.WriteString("\n")
	}

	reader := strings.NewReader(sb.String())
	participants, errors, err := ParseParticipants(reader)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	if len(participants) != 1000 {
		t.Errorf("Expected 1000 participants, got %d", len(participants))
	}

	if len(errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(errors))
	}

	t.Log("✅ Successfully parsed 1000 participants")
}

func TestValidateCSVFormat_Valid(t *testing.T) {
	csvData := `name,email,phone
John Doe,john@example.com,08123456789`

	reader := strings.NewReader(csvData)
	err := ValidateCSVFormat(reader)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateCSVFormat_Empty(t *testing.T) {
	csvData := ``

	reader := strings.NewReader(csvData)
	err := ValidateCSVFormat(reader)

	if err != domain.ErrEmptyCSV {
		t.Errorf("Expected ErrEmptyCSV, got %v", err)
	}
}

func BenchmarkParseParticipants(b *testing.B) {
	csvData := `name,email,phone
John Doe,john@example.com,08123456789
Jane Smith,jane@example.com,08987654321
Bob Johnson,bob@example.com,08111222333`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reader := strings.NewReader(csvData)
		ParseParticipants(reader)
	}
}
