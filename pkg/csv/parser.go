package csv

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/fzndps/eventcheck/internal/domain"
)

func ParseParticipants(reader io.Reader) ([]*domain.Participant, []string, error) {
	csvReader := csv.NewReader(reader)

	// Mmembaca header row
	headers, err := csvReader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	// Validate headers
	headerMap := make(map[string]int)
	for i, header := range headers {
		// Trim shitespace dan lowercase
		header = strings.TrimSpace(strings.ToLower(header))
		headerMap[header] = i
	}

	// Cek kolom yang dibutuhkan
	requiredColumns := []string{"name", "email", "phone"}
	for _, col := range requiredColumns {
		if _, exists := headerMap[col]; !exists {
			return nil, nil, fmt.Errorf("column '%s' not found in CSV", col)
		}
	}

	// Baca data rows
	var participants []*domain.Participant
	var errors []string
	rowNumber := 1 // mulai dari 1 setelah header

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			errors = append(errors, fmt.Sprintf("Rows %d: error failed read data - %v", rowNumber, err))
			rowNumber++
			continue
		}

		rowNumber++

		// ekstrak data dari row
		name := strings.TrimSpace(record[headerMap["name"]])
		email := strings.TrimSpace(record[headerMap["email"]])
		phone := strings.TrimSpace(record[headerMap["phone"]])

		// Validate data
		if name == "" {
			errors = append(errors, fmt.Sprintf("Rows %d: name empty", rowNumber))
			continue
		}

		if email == "" {
			errors = append(errors, fmt.Sprintf("Rows %d: email empty", rowNumber))
			continue
		}

		if phone == "" {
			errors = append(errors, fmt.Sprintf("Rows %d: phone empty", rowNumber))
			continue
		}

		// Basic email validasi
		if !strings.Contains(email, "@") {
			errors = append(errors, fmt.Sprintf("Rows %d: format email invalid", rowNumber))
			continue
		}

		participant := &domain.Participant{
			Name:  name,
			Email: email,
			Phone: phone,
		}

		participants = append(participants, participant)
	}

	if len(participants) == 0 && len(errors) == 0 {
		return nil, nil, domain.ErrEmptyCSV
	}

	return participants, errors, nil
}

// Melakukan validasi format csv sebelum di parse
func ValidateCSVFormat(reder io.Reader) error {
	csvReader := csv.NewReader(reder)

	_, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			return domain.ErrEmptyCSV
		}

		return domain.ErrInvalidCSVFormat
	}

	return nil
}
