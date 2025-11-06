package domain

import (
	"fmt"
	"testing"
)

func TestCalculatePrice(t *testing.T) {
	tests := []struct {
		name              string
		participantCount  int
		expectedPrice     int
		expectedPerPerson int
	}{
		{
			name:              "Tier 1: 1-50 participants",
			participantCount:  30,
			expectedPrice:     150000, // 30 * 5000
			expectedPerPerson: 5000,
		},
		{
			name:              "Tier 1 boundary: 50 participants",
			participantCount:  50,
			expectedPrice:     250000, // 50 * 5000
			expectedPerPerson: 5000,
		},
		{
			name:              "Tier 2: 51-100 participants",
			participantCount:  75,
			expectedPrice:     337500, // 75 * 4500
			expectedPerPerson: 4500,
		},
		{
			name:              "Tier 2 boundary: 100 participants",
			participantCount:  100,
			expectedPrice:     450000, // 100 * 4500
			expectedPerPerson: 4500,
		},
		{
			name:              "Tier 3: 101-500 participants",
			participantCount:  250,
			expectedPrice:     1000000, // 250 * 4000
			expectedPerPerson: 4000,
		},
		{
			name:              "Tier 3 boundary: 500 participants",
			participantCount:  500,
			expectedPrice:     2000000, // 500 * 4000
			expectedPerPerson: 4000,
		},
		{
			name:              "Tier 4: 500+ participants",
			participantCount:  1000,
			expectedPrice:     3500000, // 1000 * 3500
			expectedPerPerson: 3500,
		},
		{
			name:              "Minimum: 1 participant",
			participantCount:  1,
			expectedPrice:     5000, // 1 * 5000
			expectedPerPerson: 5000,
		},
		{
			name:              "Large event: 5000 participants",
			participantCount:  5000,
			expectedPrice:     17500000, // 5000 * 3500
			expectedPerPerson: 3500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePrice(tt.participantCount)
			fmt.Println("result calculate price test: ", result)

			if result != tt.expectedPrice {
				t.Errorf("CalculatePrice(%d) = %d, want %d",
					tt.participantCount, result, tt.expectedPrice)
			}

			// Verify price per person
			pricePerPerson := result / tt.participantCount
			fmt.Println("result price per persontest: ", pricePerPerson)
			if pricePerPerson != tt.expectedPerPerson {
				t.Errorf("Price per person = %d, want %d",
					pricePerPerson, tt.expectedPerPerson)
			}
		})
	}
}

func TestCalculatePrice_TierTransitions(t *testing.T) {
	// Test bahwa tier transitions konsisten

	// Tier 1 to Tier 2
	price50 := CalculatePrice(50)
	price51 := CalculatePrice(51)

	if price51 >= price50 {
		t.Error("Price for 51 participants should be higher than 50")
	}

	// Tier 2 to Tier 3
	price100 := CalculatePrice(100)
	price101 := CalculatePrice(101)

	if price101 >= price100 {
		t.Error("Price for 101 participants should be higher than 100")
	}

	// Tier 3 to Tier 4
	price500 := CalculatePrice(500)
	price501 := CalculatePrice(501)

	if price501 >= price500 {
		t.Error("Price for 501 participants should be higher than 500")
	}
}

func TestCalculatePrice_Discount(t *testing.T) {
	// Verify that per-person price decreases as count increases

	price30 := CalculatePrice(30)     // Tier 1: 5000/person
	price75 := CalculatePrice(75)     // Tier 2: 4500/person
	price250 := CalculatePrice(250)   // Tier 3: 4000/person
	price1000 := CalculatePrice(1000) // Tier 4: 3500/person

	pricePerPerson30 := price30 / 30
	pricePerPerson75 := price75 / 75
	pricePerPerson250 := price250 / 250
	pricePerPerson1000 := price1000 / 1000

	if pricePerPerson30 <= pricePerPerson75 {
		t.Error("Tier 1 should be more expensive per person than Tier 2")
	}
	if pricePerPerson75 <= pricePerPerson250 {
		t.Error("Tier 2 should be more expensive per person than Tier 3")
	}
	if pricePerPerson250 <= pricePerPerson1000 {
		t.Error("Tier 3 should be more expensive per person than Tier 4")
	}
}

func BenchmarkCalculatePrice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculatePrice(100)
	}
}
